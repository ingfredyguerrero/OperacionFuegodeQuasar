// Package main registra los router del servicio API REST
package main

import (
  "bytes"  
  "encoding/json"
  "fmt"
  "io/ioutil"
  "html/template"
  "log"
  "math"
  "net/http"  
  "os"

  "github.com/Microsoft/ApplicationInsights-Go/appinsights"
  "github.com/Knetic/govaluate"
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
	Message         string
	Language        string
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
	request := appinsights.NewRequestTelemetry("GET", "https://myapp.azurewebsites.net/", 1 , "Success")
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
		return
	}
	json.Unmarshal(reqBody, &req)
	for _, satel := range req.Satellites {
		writeRequest(satel.Name, satel.Distance, satel.Message)
	}
	reqUpload := cargarRequest()
	a, b := GetLocation(reqUpload)
	pos := position{X: a, Y: b}
	mess := GetMessage(reqUpload)
	resp := response{Position: pos, Message: mess}
	w.Header().Set("Content-Type", "application/json")
	if a == 9999999999 || b == 9999999999 || mess == "" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(resp)

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
	name, err := vars["satellite_name"]
	if err == false {
		fmt.Fprintln(w, "Datos Incorrectos", name)
		return
	}
	var req requestSingleSatellite
	reqBody, err2 := ioutil.ReadAll(r.Body)
	if err2 != nil {
		fmt.Fprintln(w, "Datos Incorrectos", reqBody)
		return
	}
	json.Unmarshal(reqBody, &req)
	writeRequest(name, req.Distance, req.Message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// writeRequest metodo persistir en en archivo json informacion de las naves.
func writeRequest(name string, distance float64, message []string) {
	content, err := ioutil.ReadFile("./satellites.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	req := new(request)
	json.Unmarshal([]byte(content), req)
	encontro := false

	for i := 0; i < len(req.Satellites); i++ {
		if name == req.Satellites[i].Name {
			req.Satellites[i].Distance = distance
			req.Satellites[i].Message = message

			encontro = true
			break
		}
	}

	if encontro == false {
		newSatellite := satellites{Name:name, Distance:distance, Message:message}
		req.Satellites = append(req.Satellites, newSatellite)
	}

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
	var reqUpload = cargarRequest()
	var a, b = GetLocation(reqUpload)
	var pos position
	pos.X = a
	pos.Y = b
	mess := GetMessage(reqUpload)
	var resp response
	resp.Position = pos
	resp.Message = mess
	w.Header().Set("Content-Type", "application/json")
	if a == 9999999999 || b == 9999999999 || mess == "" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(resp)
}

// cargarRequest metodo para cargar la informacion persistida de la nave y distancia de los satelites.
func cargarRequest() *request {
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
		Message: "Success!",
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
	sum := 0
	expression, err := govaluate.NewEvaluableExpression("((x2 - x1) ** 2)+((y2 - y1) ** 2)")
	if err != nil {
		return
	}
	satt := traerSatelites()

	for x := -110.0; x < 200; x = x + 0.1 {
		for y := -110.0; y < 200; y = y + 0.1 {
			for _, satelDis := range distances.Satellites {
				for _, satel := range satt.Satelites {
					if satel.Nombre == satelDis.Name {
						parameters := make(map[string]interface{}, 8)
						parameters["x1"] = x
						parameters["y1"] = y
						parameters["x2"] = roundTo(satel.X, 2)
						parameters["y2"] = roundTo(satel.Y, 2)
						result, err := expression.Evaluate(parameters)

						if err != nil {
							break
						}
						raiz := math.Sqrt(result.(float64))
						if roundTo(satelDis.Distance, 2) == roundTo(raiz, 2) {
							sum++
						}
						break
					}

				}

			}
			if sum >= 2 {
				return roundTo(x, 2), roundTo(y, 2)
			}
			sum = 0

		}
	}
	return w, z
}

// Metodo para redondear a dos decimales
func roundTo(n float64, decimals uint32) float64 {
	return math.Round(n*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
}

// GetMessage realiza la busqueda del mensage de le nave a la que se interceptaron los mensajes
// input: el mensaje tal cual es recibido en cada satélite
// output: el mensaje tal cual lo genera el emisor del mensaje
func GetMessage(messages *request) (msg string) {
	var respuesta []string
	retorna := ""
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