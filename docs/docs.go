// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "mymail",
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
        "/streamer": {
            "get": {
                "description": "Request a random streamer from a country",
                "produces": [
                    "application/json"
                ],
                "summary": "/streamer?lang=fr\u0026minviewers=1\u0026maxviewers=10",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Language of the streamer, should be an ISO code like fr,de,it. Default to fr",
                        "name": "lang",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Minimum number of viewers you want the streamer to have. Default to 1",
                        "name": "minviewers",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Maximum number of viewers you want the streamer to have. Default to 10",
                        "name": "maxviewers",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Streamer"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.Streamer": {
            "type": "object",
            "properties": {
                "gameName": {
                    "type": "string"
                },
                "userName": {
                    "type": "string"
                },
                "viewerCount": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "twitch-highlight-streamer",
	Description:      "This is an api to retrieve a random streamer on Twtich.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}