// Package main registra los router del servicio API REST
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/Knetic/govaluate"
	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"github.com/gorilla/mux"
	//httpSwagger "github.com/swaggo/http-swagger"
)

// satelites representa los satelites existentes con sus coordenadas
type satelites struct {
	Nombre string  `json:"nombre"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
}

// response representa la respuesta del mensaje con sus coordenadas
type response struct {
	Position position `json:"position"`
	Message  string   `json:"message"`
}

// response representa la respuesta del mensaje con sus coordenadas
type responseError struct {
	Message string `json:"message"`
}

// position representa las coordenadas
type position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// satellites representa la informacion de la nave y distancia de los satelites
type satellites struct {
	Name     string   `json:"name"`
	Distance float64  `json:"distance"`
	Message  []string `json:"message"`
}

// request representa el objeto de entrada de la api
type request struct {
	Satellites []satellites `json:"satellites"`
}

// configuration representa la configuracion de los satelites actuales
type configuration struct {
	Satelites []satelites `json:"satelites"`
}

// PageVars representa informacion de la pagina index
type PageVars struct {
	Message  string
	Language string
}

// requestSingleSatellite representa la entrada de la api con informacion de solo un satelite
type requestSingleSatellite struct {
	Distance float64  `json:"distance"`
	Message  []string `json:"message"`
}

// main inicializa el servicio con el puerto y los router de cada api
// @title Satellite API
// @version 1.0
// @description This is a service for calculate the Lotation of consumers
// @termsOfService http://swagger.io/terms/
// @contact.name API Satellite
// @contact.email ingfredyguerrero@hotmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:80
// @BasePath /
func main() {
	r := mux.NewRouter().StrictSlash(true)
	client := appinsights.NewTelemetryClient(os.Getenv("APPINSIGHTS_INSTRUMENTATIONKEY"))
	request := appinsights.NewRequestTelemetry("GET", "https://myapp.azurewebsites.net/", 1, "Success")
	client.Track(request)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("fonts"))))
	http.HandleFunc("/", Home)
	r.HandleFunc("/topsecret", topSecret).Methods("POST")
	r.HandleFunc("/topsecret_split/{satellite_name}", topSecretSplit).Methods("POST")
	r.HandleFunc("/satellite/topsecret_split", getSatellites).Methods("GET")
	//r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	log.Fatal(http.ListenAndServe(getPort(), r))
}

// topSecret godoc
// @Summary metodo para recibir las cordenadas de la nave con los mensajes y realizar el proceso de ubicacion y codificacion del mensaje
// @Description metodo para recibir las cordenadas de la nave con los mensajes y realizar el proceso de ubicacion y codificacion del mensaje
// @Tags intercept
// @Accept  json
// @Produce  json
// @Success 200
// @Router /topSecret [post]
func topSecret(w http.ResponseWriter, r *http.Request) {
	req := request{}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, "Datos Incorrectos", reqBody)
		errorMessage := responseError{Message: "Datos Incorrectos"}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// transforma los datos en la entidad request
	json.Unmarshal(reqBody, &req)
	//recorre los satelites enviados por el cosumidor
	for _, satel := range req.Satellites {
		// Persiste los satelites en el archivo satelites.json
		writeRequest(satel.Name, satel.Distance, satel.Message)
	}
	// carga los satelites persistidos en el archivo satelites.json
	reqUpload := readRequest()
	// obtiene la ubicacion de la nave que envio los llamados de auxilio.
	a, b := GetLocation(reqUpload)
	pos := position{X: a, Y: b}
	// obtiene los mensajes de la nave que envio los llamados de auxilio.
	mess := GetMessage(reqUpload)
	resp := response{Position: pos, Message: mess}
	w.Header().Set("Content-Type", "application/json")
	//Valida la informacion calculada
	if a == 9999999999 || b == 9999999999 || mess == "" {
		w.WriteHeader(http.StatusNotFound)
		errorMessage := responseError{Message: "No se pudo calcular la informacion de la nave"}
		json.NewEncoder(w).Encode(errorMessage)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

// topSecretSplit godoc
// @Summary metodo para almacenar informacion de la nave en relacion a solo un satelite
// @Description metodo para almacenar informacion de la nave en relacion a solo un satelite
// @Tags intercept
// @Accept  json
// @Produce  json
// @Param RequestSingleSatellite body requestSingleSatellite true "Distance of Satellite"
// @Param satellite_name path string true "name of Satellite"
// @Success 200
// @Router /topsecret_split/{satellite_name} [post]
func topSecretSplit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//valida infornacion del path
	name, err := vars["satellite_name"]
	if err == false {
		errorMessage := responseError{Message: "Datos Incorrectos"}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var req requestSingleSatellite
	//valida infornacion del body
	reqBody, err2 := ioutil.ReadAll(r.Body)
	if err2 != nil {
		fmt.Fprintln(w, "Datos Incorrectos", reqBody)
		errorMessage := responseError{Message: "Datos Incorrectos"}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	json.Unmarshal(reqBody, &req)
	//Persiste la informacion del satelite en el archivo satelites.json
	writeRequest(name, req.Distance, req.Message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// writeRequest metodo persistir en en archivo satelites.json con la informacion de las naves.
func writeRequest(name string, distance float64, message []string) {
	// lee el archivo actual
	content, err := ioutil.ReadFile("./satellites.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	req := new(request)
	json.Unmarshal([]byte(content), req)
	encontro := false
	// busca si existe el satelite
	for i := 0; i < len(req.Satellites); i++ {
		// si encuentra el satelite lo actualiza.
		if name == req.Satellites[i].Name {
			req.Satellites[i].Distance = distance
			req.Satellites[i].Message = message

			encontro = true
			break
		}
	}
	// NO si encuentra el satelite lo Crea.
	if encontro == false {
		newSatellite := satellites{Name: name, Distance: distance, Message: message}
		req.Satellites = append(req.Satellites, newSatellite)
	}

	//Actualiza el archivo con el nuevo listado de satelites
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(req)
	ioutil.WriteFile("./satellites.json", reqBodyBytes.Bytes(), os.ModePerm)

}

// getSatellites godoc
// @Summary metodo para obtener ubicacion y mensajes decodificados de la nave.
// @Description metodo para obtener ubicacion y mensajes decodificados de la nave.
// @Tags intercept
// @Accept  json
// @Produce  json
// @Success 200
// @Router /topsecret_split [get]
func getSatellites(w http.ResponseWriter, r *http.Request) {
	var reqUpload = readRequest()
	// consulta la ubicacion de la nave de la informacion previamente recibida
	var a, b = GetLocation(reqUpload)
	var pos position
	pos.X = a
	pos.Y = b
	// consulta llos mensajes de la informacion previamente recibida
	mess := GetMessage(reqUpload)
	var resp response
	resp.Position = pos
	resp.Message = mess
	w.Header().Set("Content-Type", "application/json")
	//Valida la informacion calculada
	if a == 9999999999 || b == 9999999999 || mess == "" {
		w.WriteHeader(http.StatusNotFound)
		errorMessage := responseError{Message: "No se pudo calcular la informacion de la nave"}
		json.NewEncoder(w).Encode(errorMessage)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

// readRequest metodo para cargar la informacion persistida de la nave y distancia de los satelites.
func readRequest() *request {
	// lee y retorna la informacion del archivo de satelites.json
	content, err := ioutil.ReadFile("./satellites.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	req := new(request)
	json.Unmarshal([]byte(content), req)
	return req
}

// getPort metodo para obtener el puerto de publicacion del servicio
func getPort() string {
	p := os.Getenv("HTTP_PLATFORM_PORT")
	if p != "" {
		return ":" + p
	}
	return ":80"
}

// render metodo autogenerado para index
func render(w http.ResponseWriter, tmpl string, pageVars PageVars) {

	tmpl = fmt.Sprintf("views/%s", tmpl)
	t, err := template.ParseFiles(tmpl)

	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}

	err = t.Execute(w, pageVars) //execute the template and pass in the variables to fill the gaps

	if err != nil { // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

// Home metodo para cargar pagina Index
func Home(w http.ResponseWriter, req *http.Request) {
	pageVars := PageVars{
		Message:  "Success!",
		Language: "Go Lang",
	}
	render(w, "index.html", pageVars)
}

// GetLocation realiza la busqueda de las coordenadas de le nave a la que se interceptaron los mensajes
// input: distancia al emisor tal cual se recibe en cada satélite
// output: las coordenadas ‘x’ e ‘y’ del emisor del mensaje
func GetLocation(distances *request) (w, z float64) {
	w = 9999999999
	z = 9999999999
	for _, satelDis1 := range distances.Satellites {
		// inicia el recorrido de los satelites, obtiene la ubicacion del satelite por nombre               fmt.Println("entro satelDis1")
		X1, Y1 := findSatelite(satelDis1.Name)

		for _, satelDis2 := range distances.Satellites {
			// valida que no se realicen calculos para el mismo satelite           fmt.Println("entro satelDis2")
			if satelDis1.Name != satelDis2.Name {
				// obtiene la ubicacion del segundo satelite a evaluar por nombre
				X2, Y2 := findSatelite(satelDis2.Name)
				// calcula la variable X positiva usando la formula de despeje de la ecuacion de distancia
				Xpositive, err := calculateX(satelDis1.Distance, X1, Y1, satelDis2.Distance, X2, Y2, true)
				if err != nil {
					break
				}
				// calcula la variable Y positiva  usando la variable X que calculo de la anterior ecuacion usando la formula de despeje de la ecuacion de distancia
				Yposivite1, err := calculateY(satelDis1.Distance, X1, Y1, Xpositive)
				if err != nil {
					break
				}
				// validacion de resultado X positivo con primera Y usando la formula inicial para identificar si la distancia da igual con los valores calculados
				Validate1, err := ValidatePoints(satelDis1.Distance, X1, Y1, satelDis2.Distance, X2, Y2, Xpositive, Yposivite1)
				if err != nil {
					break
				}
				//fmt.Println("X: ", roundTo(Xpositive, 0), " Y: ", roundTo(Yposivite1, 0), " x1: ", X1, " Y1: ", Y1, " Dis1: ", roundTo(satelDis1.Distance, 2), " x2: ", X2, " Y2: ", Y2, " Dis2: ", roundTo(satelDis2.Distance, 2))
				// si la validacion es exitosa encontro los valores y los retorna
				if Validate1 {
					return roundTo(Xpositive, 0), roundTo(Yposivite1, 0)
				}
				// calcula la variable Y segunda formula usando la variable X que calculo de la anterior ecuacion usando la formula de despeje de la ecuacion de distancia
				Yposivite2, err := calculateY(satelDis2.Distance, X2, Y2, Xpositive)
				if err != nil {
					break
				}
				// validacion de resultado X positivo con segunda Y usando la formula inicial para identificar si la distancia da igual con los valores calculados
				Validate2, err := ValidatePoints(satelDis1.Distance, X1, Y1, satelDis2.Distance, X2, Y2, Xpositive, Yposivite2)
				if err != nil {
					break
				}
				//fmt.Println("X: ", roundTo(Xpositive, 0), " Y: ", roundTo(Yposivite2, 0), " x1: ", X1, " Y1: ", Y1, " Dis1: ", roundTo(satelDis1.Distance, 2), " x2: ", X2, " Y2: ", Y2, " Dis2: ", roundTo(satelDis2.Distance, 2))
				// si la validacion es exitosa encontro los valores y los retorna
				if Validate2 {
					return roundTo(Xpositive, 0), roundTo(Yposivite2, 0)
				}
				// calcula la variable X negativa usando la formula de despeje de la ecuacion de distancia
				Xnegative, err := calculateX(satelDis1.Distance, X1, Y1, satelDis2.Distance, X2, Y2, false)
				if err != nil {
					break
				}
				// calcula la variable Y negativa usando la variable X que calculo de la anterior ecuacion usando la formula de despeje de la ecuacion de distancia
				Ynegative1, err := calculateY(satelDis1.Distance, X1, Y1, Xnegative)
				if err != nil {
					break
				}
				// validacion de resultado X negativo con primer Y usando la formula inicial para identificar si la distancia da igual con los valores calculados
				Validate3, err := ValidatePoints(satelDis1.Distance, X1, Y1, satelDis2.Distance, X2, Y2, Xnegative, Ynegative1)
				if err != nil {
					break
				}
				//fmt.Println("X: ", roundTo(Xnegative, 0), " Y: ", roundTo(Ynegative1, 0), " x1: ", X1, " Y1: ", Y1, " Dis1: ", roundTo(satelDis1.Distance, 2), " x2: ", X2, " Y2: ", Y2, " Dis2: ", roundTo(satelDis2.Distance, 2))
				// si la validacion es exitosa encontro los valores y los retorna
				if Validate3 {
					return roundTo(Xnegative, 0), roundTo(Ynegative1, 0)
				}
				// calcula la variable Y segunda formula usando la variable X que calculo de la anterior ecuacion usando la formula de despeje de la ecuacion de distancia
				Ynegative2, err := calculateY(satelDis2.Distance, X2, Y2, Xnegative)
				if err != nil {
					break
				}
				// validacion de resultado X negativo con la segunda Y usando la formula inicial para identificar si la distancia da igual con los valores calculados
				Validate4, err := ValidatePoints(satelDis1.Distance, X1, Y1, satelDis2.Distance, X2, Y2, Xnegative, Ynegative2)
				if err != nil {
					break
				}
				//fmt.Println("X: ", roundTo(Xnegative, 0), " Y: ", roundTo(Ynegative2, 0), " x1: ", X1, " Y1: ", Y1, " Dis1: ", roundTo(satelDis1.Distance, 2), " x2: ", X2, " Y2: ", Y2, " Dis2: ", roundTo(satelDis2.Distance, 2))
				// si la validacion es exitosa encontro los valores y los retorna
				if Validate4 {
					return roundTo(Xnegative, 0), roundTo(Ynegative2, 0)
				}
			}
		}
	}

	/// codigo comentariado de calcular con el recorrido de la matriz en el plano carteciano los valores

	// expression, err := govaluate.NewEvaluableExpression("((x2 - x1) ** 2)+((y2 - y1) ** 2)")
	// if err != nil {
	// 	return
	// }

	// for x := -110.0; x < 200; x = x + 0.1 {
	// 	for y := -110.0; y < 200; y = y + 0.1 {
	// 		for _, satelDis := range distances.Satellites {
	// 			for _, satel := range satt.Satelites {
	// 				if satel.Nombre == satelDis.Name {
	// 					parameters := make(map[string]interface{}, 8)
	// 					parameters["x1"] = x
	// 					parameters["y1"] = y
	// 					parameters["x2"] = roundTo(satel.X, 2)
	// 					parameters["y2"] = roundTo(satel.Y, 2)
	// 					result, err := expression.Evaluate(parameters)

	// 					if err != nil {
	// 						break
	// 					}
	// 					raiz := math.Sqrt(result.(float64))
	// 					if roundTo(satelDis.Distance, 2) == roundTo(raiz, 2) {
	// 						sum++
	// 					}
	// 					break
	// 				}

	// 			}

	// 		}
	// 		if sum >= 2 {
	// 			return roundTo(x, 2), roundTo(y, 2)
	// 		}
	// 		sum = 0

	// 	}
	// }

	// si no enncuentra lo valores retorna los valores por defecto para ser evaluados
	return w, z
}

// Metodo para redondear a dos decimales
func roundTo(n float64, decimals uint32) float64 {
	return math.Round(n*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
}

// ValidatePoints Metodo validar dos puntos
func ValidatePoints(distance1 float64, x1 float64, y1 float64, distance2 float64, x2 float64, y2 float64, xCalculate float64, yCalculate float64) (bool, error) {
	// formula ecuacion para calcular la primer distancia
	expressionVerificacion, err := govaluate.NewEvaluableExpression("((x2 - x1) ** 2)+((y2 - y1) ** 2)")
	if err != nil {
		fmt.Println("entro error expressionVerificacion ", err)
		return false, err
	}
	// parametros de la formula para calcular la primer distancia
	parametersVerificacionDistancia1 := make(map[string]interface{}, 8)
	parametersVerificacionDistancia1["x1"] = roundTo(xCalculate, 0)
	parametersVerificacionDistancia1["y1"] = roundTo(yCalculate, 0)
	parametersVerificacionDistancia1["x2"] = x1
	parametersVerificacionDistancia1["y2"] = y1
	// calcula el resultado de la distancia
	resultDistancia1, err := expressionVerificacion.Evaluate(parametersVerificacionDistancia1)

	if err != nil {
		return false, err
	}
	// parametros de la formula para calcular la primer distancia
	parametersVerificacionDistancia2 := make(map[string]interface{}, 8)
	parametersVerificacionDistancia2["x1"] = roundTo(xCalculate, 0)
	parametersVerificacionDistancia2["y1"] = roundTo(yCalculate, 0)
	parametersVerificacionDistancia2["x2"] = x2
	parametersVerificacionDistancia2["y2"] = y2
	// calcula el resultado de la distancia
	resultDistancia2, err := expressionVerificacion.Evaluate(parametersVerificacionDistancia2)

	if err != nil {
		return false, err
	}
	fmt.Println(distance1, " Dis1: ", roundTo(math.Sqrt(resultDistancia1.(float64)), 0), " ", distance2, " Dis2: ", roundTo(math.Sqrt(resultDistancia2.(float64)), 0))
	// valida si las distancias encontradas corresponden a las de los satelites
	if roundTo(distance1, 0) == roundTo(math.Sqrt(resultDistancia1.(float64)), 0) && roundTo(distance2, 0) == roundTo(math.Sqrt(resultDistancia2.(float64)), 0) {
		return true, nil
	}
	return false, nil
}

// findSatelite busca la ubicacion de un satelite por nombre
func findSatelite(name string) (x float64, y float64) {
	satt := traerSatelites()
	for _, satelLocation2 := range satt.Satelites {
		if name == satelLocation2.Nombre {
			return satelLocation2.X, satelLocation2.Y
		}
	}
	return 0, 0
}

// calculateY Metodo para calcular Y
func calculateY(distance float64, x float64, y float64, resultX float64) (float64, error) {
	// formula sin raiz praa calular Y
	expressionY, err := govaluate.NewEvaluableExpression("(w - r)")
	if err != nil {
		fmt.Println("entro error expressionY ", err)
		return 0, nil
	}
	// formula con raiz praa calular Y
	expressionYRaiz, err := govaluate.NewEvaluableExpression("( (s**2) - ((t - x)**2) )")
	if err != nil {
		fmt.Println("entro error expressionYRaiz ", err)
		return 0, nil
	}
	// parametros para la formula con raiz
	parametersYRaiz := make(map[string]interface{}, 8)
	parametersYRaiz["s"] = roundTo(distance, 2)
	parametersYRaiz["t"] = x
	parametersYRaiz["x"] = roundTo(resultX, 2)
	// calcula Y parcial
	resultYRaiz, err := expressionYRaiz.Evaluate(parametersYRaiz)
	if err != nil {
		return 0, nil
	}
	// parametros para la formula sin raiz
	parametersY := make(map[string]interface{}, 8)
	parametersY["r"] = roundTo(math.Sqrt(resultYRaiz.(float64)), 2)
	parametersY["w"] = y
	// calcula Y final
	resultY, err := expressionY.Evaluate(parametersY)
	if err != nil {
		return 0, nil
	}
	return resultY.(float64), nil
}

// calculateY Metodo para calcular Y
func calculateX(distance1 float64, x1 float64, y1 float64, distance2 float64, x2 float64, y2 float64, posivite bool) (float64, error) {
	result := 0.0
	// formula sin raiz praa calular X
	expressionRaiz, err := govaluate.NewEvaluableExpression("(((-4*(h**3)) + (4*(h**2)*t) + (4*h*(k**2)) - (4*h*(p**2)) + (8*h*p*w) - (4*h*(s**2)) + (4*h*(t**2)) - (4*h*(w**2)) - (4*(k**2)*t) - (4*(p**2)*t) + (8*p*t*w) + (4*(s**2)*t) - (4*(t**3)) - (4*t*(w**2)))**2 - 4*((4*(h**2)) - (8*h*t) + (4*(p**2)) - (8*p*w) + (4*(t**2)) + (4*(w**2)))*((h**4) - (2*(h**2)*(k**2)) + (2*(h**2)*(p**2)) - (4*(h**2)*p*w) + (2*(h**2)*(s**2)) - (2*(h**2)*(t**2)) + (2*(h**2)*(w**2)) + (k**4) - (2*(k**2)*(p**2)) + (4*(k**2)*p*w) - (2*(k**2)*(s**2)) + (2*(k**2)*(t**2)) - (2*(k**2)*(w**2)) + (p**4) - (4*(p**3)*w) - (2*(p**2)*(s**2)) + (2*(p**2)*(t**2)) + (6*(p**2)*(w**2)) + (4*p*(s**2)*w) - (4*p*(t**2)*w) - (4*p*(w**3)) + (s**4) - (2*(s**2)*(t**2)) - (2*(s**2)*(w**2)) + (t**4) + (2*(t**2)*(w**2)) + (w**4)))")
	if err != nil {
		fmt.Println("entro error expressionRaiz ", err)
		return 0, err
	}
	// formula completa variable positiva para calular X
	expressionXPos, err := govaluate.NewEvaluableExpression("((4*(h**3)) - (4*(h**2)*t) + r - (4*h*(k**2)) + (4*h*(p**2)) - (8*h*p*w) + (4*h*(s**2)) - (4*h*(t**2)) + (4*h*(w**2)) + (4*(k**2)*t) + (4*(p**2)*t) - (8*p*t*w) - (4*(s**2)*t) + (4*(t**3)) + (4*t*(w**2)))/(2*((4*(h**2)) - (8*h*t) + (4*(p**2)) - (8*p*w) + (4*(t**2)) + (4*(w**2))))")
	if err != nil {
		fmt.Println("entro error expressionXPos ", err)
		return 0, err
	}
	// formula completa variable negativa para calular X
	expressionXNeg, err := govaluate.NewEvaluableExpression("((4*(h**3)) - (4*(h**2)*t) - r - (4*h*(k**2)) + (4*h*(p**2)) - (8*h*p*w) + (4*h*(s**2)) - (4*h*(t**2)) + (4*h*(w**2)) + (4*(k**2)*t) + (4*(p**2)*t) - (8*p*t*w) - (4*(s**2)*t) + (4*(t**3)) + (4*t*(w**2)))/(2*((4*(h**2)) - (8*h*t) + (4*(p**2)) - (8*p*w) + (4*(t**2)) + (4*(w**2))))")
	if err != nil {
		fmt.Println("entro error expressionXNeg ", err)
		return 0, err
	}
	// parametros para la formula sin raiz
	parametersRaiz := make(map[string]interface{}, 8)
	parametersRaiz["s"] = distance1
	parametersRaiz["t"] = x1
	parametersRaiz["w"] = y1
	parametersRaiz["k"] = distance2
	parametersRaiz["h"] = x2
	parametersRaiz["p"] = y2
	// calculo de la variable sin razin
	resultRaiz, err := expressionRaiz.Evaluate(parametersRaiz)

	if err != nil {
		return 0, err
	}
	// parametros para la formula completa X
	parametersXPos := make(map[string]interface{}, 8)
	parametersXPos["s"] = distance1
	parametersXPos["t"] = x1
	parametersXPos["w"] = y1
	parametersXPos["k"] = distance2
	parametersXPos["h"] = x2
	parametersXPos["p"] = y2
	parametersXPos["r"] = roundTo(math.Sqrt(resultRaiz.(float64)), 2)
	//Valida si debe aplicar la formula para calculo del valor positivo (true) o negativo (false)
	if posivite == true {
		//calcula valor positivo
		resultX, err := expressionXPos.Evaluate(parametersXPos)
		if err != nil {
			return 0, err
		}
		result = resultX.(float64)
	} else {
		//calcula valor negativo
		resultX, err := expressionXNeg.Evaluate(parametersXPos)
		if err != nil {
			return 0, err
		}
		result = resultX.(float64)
	}

	return result, nil
}

// GetMessage realiza la busqueda del mensage de le nave a la que se interceptaron los mensajes
// input: el mensaje tal cual es recibido en cada satélite
// output: el mensaje tal cual lo genera el emisor del mensaje
func GetMessage(messages *request) (msg string) {
	var respuesta []string
	retorna := ""
	//recorre los satelites y evalua sus mensajes para determinar si coinciden o no de acuerdo a su posicion
	for i := 0; i < len(messages.Satellites); i++ {
		for j := 0; j < len(messages.Satellites[i].Message); j++ {
			contarRespuesta := len(respuesta)
			if contarRespuesta <= j {
				respuesta = append(respuesta, messages.Satellites[i].Message[j])
			} else if respuesta[j] == "" {
				respuesta[j] = messages.Satellites[i].Message[j]
			} else if messages.Satellites[i].Message[j] != "" && respuesta[j] != messages.Satellites[i].Message[j] {
				respuesta[j] = "COMODIN"
			}
		}
	}
	for _, mgsItem := range respuesta {
		if retorna == "" && mgsItem != "COMODIN" {
			retorna = mgsItem
		} else if mgsItem != "COMODIN" {
			retorna = retorna + " " + mgsItem
		}
	}
	return retorna
}

// traerSatelites Metodo para traer la configuracion en json de los satelites y cordenadas actuales
// si se requiere aumentar los satelites solo se debe modificar el archivo config.json y tomara los nuevos satelites para los calculos
func traerSatelites() *configuration {
	content, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	conf := new(configuration)
	json.Unmarshal([]byte(content), conf)
	return conf
}
