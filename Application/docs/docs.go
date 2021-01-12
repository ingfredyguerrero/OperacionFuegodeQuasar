// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Satellite",
            "email": "ingfredyguerrero@hotmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/topSecret": {
            "post": {
                "description": "set distance of satellites and get location of consumers",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "intercept"
                ],
                "summary": "set distance of satellites and get location of consumers",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/topsecret_split": {
            "get": {
                "description": "Get location of consumers",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "intercept"
                ],
                "summary": "Get location of consumers",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/topsecret_split/{satellite_name}": {
            "post": {
                "description": "set location of satellite",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "intercept"
                ],
                "summary": "set location of satellite",
                "parameters": [
                    {
                        "description": "Distance of Satellite",
                        "name": "RequestSingleSatellite",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.requestSingleSatellite"
                        }
                    },
                    {
                        "type": "string",
                        "description": "name of Satellite",
                        "name": "satellite_name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        }
    },
    "definitions": {
        "main.requestSingleSatellite": {
            "type": "object",
            "properties": {
                "distance": {
                    "type": "number"
                },
                "message": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:9876",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Satellite API",
	Description: "This is a service for calculate the Lotation of consumers",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
