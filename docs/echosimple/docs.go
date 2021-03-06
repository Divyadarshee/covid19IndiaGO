// Package echosimple GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package echosimple

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
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
        "/": {
            "get": {
                "description": "get the status of server.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "Show the status of server.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/cases": {
            "get": {
                "description": "Get the covid19 cases details for the given GPS coordinates.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cases"
                ],
                "summary": "Get latest covid19 cases.",
                "parameters": [
                    {
                        "type": "number",
                        "description": "Latitude",
                        "name": "latitude",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "Longitude",
                        "name": "longitude",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "covid19-go-deploy.herokuapp.com",
	BasePath:         "/",
	Schemes:          []string{"https"},
	Title:            "Covid19 Cases and Vaccinations in India",
	Description:      "This is a covid19 cases data server which when given the GPS coordinates of a location returns\nthe cases details as in confirmed, deceased, recovered, teseted along with vaccination details as in\nsingle and double dose in coordinates provided State and in India in total",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
