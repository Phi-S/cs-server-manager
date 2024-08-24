// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/logs/{countOrSince}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "logs"
                ],
                "summary": "Gets logs",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Gets the last x logs or all logs since date",
                        "name": "countOrSince",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/logwrt.LogEntry"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/plugins": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "plugins"
                ],
                "summary": "Gets all available plugins",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.PluginResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "plugins"
                ],
                "summary": "Installs the given plugin or updates to given version",
                "parameters": [
                    {
                        "description": "The plugin and the version that should be installed",
                        "name": "plugin",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.InstallPluginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "tags": [
                    "plugins"
                ],
                "summary": "Uninstalls the currently installed plugin",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/send-command": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "server"
                ],
                "summary": "Sends and executes a game server command",
                "parameters": [
                    {
                        "description": "This command will be executed on the game server",
                        "name": "command",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.CommandRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.CommandResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/settings": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "settings"
                ],
                "summary": "Gets the current settings",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.SettingsModel"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "settings"
                ],
                "summary": "Gets the current settings",
                "parameters": [
                    {
                        "description": "The settings to update",
                        "name": "settings",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.SettingsModel"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.SettingsModel"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/start": {
            "post": {
                "description": "Starts the server with the given start parameters",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "server"
                ],
                "summary": "Starts the server",
                "parameters": [
                    {
                        "description": "You can provide no, all or only a few start parameters. The provided start parameters will overwrite the saved start parameters in the start-parameters.json file.",
                        "name": "startParameters",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.StartBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/status": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "server"
                ],
                "summary": "Get the current status of the server",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/status.InternalStatus"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/stop": {
            "post": {
                "description": "If the server is not running, returns 200 OK",
                "tags": [
                    "server"
                ],
                "summary": "Stops the server",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/update": {
            "post": {
                "tags": [
                    "update"
                ],
                "summary": "Starts server update",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/update/cancel": {
            "post": {
                "tags": [
                    "update"
                ],
                "summary": "Cancels the server update",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.CommandRequest": {
            "type": "object",
            "required": [
                "command"
            ],
            "properties": {
                "command": {
                    "type": "string"
                }
            }
        },
        "handlers.CommandResponse": {
            "type": "object",
            "properties": {
                "output": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "handlers.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "request_id": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "handlers.InstallPluginRequest": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "handlers.PluginResponse": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                },
                "versions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handlers.PluginVersionResponse"
                    }
                }
            }
        },
        "handlers.PluginVersionDependenciesResponse": {
            "type": "object",
            "properties": {
                "plugin_name": {
                    "type": "string"
                },
                "version_name": {
                    "type": "string"
                }
            }
        },
        "handlers.PluginVersionResponse": {
            "type": "object",
            "properties": {
                "dependencies": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handlers.PluginVersionDependenciesResponse"
                    }
                },
                "installed": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "handlers.SettingsModel": {
            "type": "object",
            "required": [
                "hostname",
                "max_players",
                "start_map"
            ],
            "properties": {
                "hostname": {
                    "type": "string",
                    "maxLength": 128
                },
                "max_players": {
                    "type": "integer",
                    "maximum": 128
                },
                "password": {
                    "type": "string",
                    "maxLength": 32
                },
                "start_map": {
                    "type": "string",
                    "maxLength": 32
                },
                "steam_login_token": {
                    "type": "string"
                }
            }
        },
        "handlers.StartBody": {
            "type": "object",
            "properties": {
                "hostname": {
                    "type": "string",
                    "maxLength": 128
                },
                "max_players": {
                    "type": "integer",
                    "maximum": 128
                },
                "password": {
                    "type": "string",
                    "maxLength": 32
                },
                "start_map": {
                    "type": "string",
                    "maxLength": 32
                },
                "steam_login_token": {
                    "type": "string"
                }
            }
        },
        "logwrt.LogEntry": {
            "type": "object",
            "properties": {
                "log_type": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                }
            }
        },
        "status.InternalStatus": {
            "type": "object",
            "properties": {
                "hostname": {
                    "type": "string"
                },
                "ip": {
                    "type": "string"
                },
                "map": {
                    "type": "string"
                },
                "max_player_count": {
                    "type": "integer"
                },
                "password": {
                    "type": "string"
                },
                "player_count": {
                    "type": "integer"
                },
                "port": {
                    "type": "string"
                },
                "state": {
                    "$ref": "#/definitions/status.State"
                }
            }
        },
        "status.State": {
            "type": "string",
            "enum": [
                "idle",
                "server-starting",
                "server-started",
                "server-stopping",
                "steamcmd-updating",
                "plugin-installing",
                "plugin-uninstalling"
            ],
            "x-enum-varnames": [
                "Idle",
                "ServerStarting",
                "ServerStarted",
                "ServerStopping",
                "SteamcmdUpdating",
                "PluginInstalling",
                "PluginUninstalling"
            ]
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{"http", "https"},
	Title:            "cs-server-manager API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
