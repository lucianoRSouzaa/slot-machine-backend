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
        "/machines": {
            "post": {
                "description": "Permite a criação de uma nova máquina caça-níqueis com os parâmetros especificados.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "SlotMachine"
                ],
                "summary": "Criar uma nova máquina caça-níqueis",
                "parameters": [
                    {
                        "description": "Dados da máquina caça-níqueis a ser criada",
                        "name": "createSlotMachineRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/usecase.CreateSlotMachineRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Máquina criada com sucesso",
                        "schema": {
                            "$ref": "#/definitions/usecase.CreateSlotMachineResponse"
                        }
                    },
                    "400": {
                        "description": "Payload inválido ou parâmetros inválidos",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    },
                    "409": {
                        "description": "Máquina caça-níqueis já existe",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Erro interno do servidor",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    }
                }
            }
        },
        "/play": {
            "post": {
                "description": "Permite que o jogador faça uma aposta e jogue na máquina caça-níqueis especificada.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "SlotMachine"
                ],
                "summary": "Jogar na máquina caça-níqueis",
                "parameters": [
                    {
                        "description": "Dados da jogada",
                        "name": "playRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/usecase.PlayRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Jogada realizada com sucesso",
                        "schema": {
                            "$ref": "#/definitions/usecase.PlayResponse"
                        }
                    },
                    "400": {
                        "description": "Payload inválido",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Máquina caça-níqueis não encontrada",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    },
                    "422": {
                        "description": "Saldo insuficiente",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Erro interno do servidor",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    }
                }
            }
        },
        "/players": {
            "post": {
                "description": "Permite a criação de um novo jogador com um saldo inicial.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Player"
                ],
                "summary": "Criar um novo jogador",
                "parameters": [
                    {
                        "description": "Dados do jogador a ser criado",
                        "name": "createPlayerRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/usecase.CreatePlayerRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Jogador criado com sucesso",
                        "schema": {
                            "$ref": "#/definitions/usecase.CreatePlayerResponse"
                        }
                    },
                    "400": {
                        "description": "Payload inválido",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    },
                    "409": {
                        "description": "Jogador já existe",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Erro interno do servidor",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "http.HTTPError": {
            "description": "Estrutura para representar erros na API. Contém a mensagem de erro e um código opcional. Pode ser expandida conforme necessário.",
            "type": "object",
            "properties": {
                "code": {
                    "description": "Código do erro HTTP",
                    "type": "integer"
                },
                "message": {
                    "description": "Mensagem descritiva do erro",
                    "type": "string"
                }
            }
        },
        "model.Player": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "model.SlotMachine": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "integer"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "initial_balance": {
                    "type": "integer"
                },
                "level": {
                    "type": "integer"
                },
                "multiple_gain": {
                    "type": "integer"
                },
                "permutations": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "string"
                        }
                    }
                },
                "symbols": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                }
            }
        },
        "usecase.CreatePlayerRequest": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "usecase.CreatePlayerResponse": {
            "type": "object",
            "properties": {
                "player": {
                    "$ref": "#/definitions/model.Player"
                }
            }
        },
        "usecase.CreateSlotMachineRequest": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "integer"
                },
                "description": {
                    "type": "string"
                },
                "level": {
                    "type": "integer"
                },
                "multiple_gain": {
                    "type": "integer"
                }
            }
        },
        "usecase.CreateSlotMachineResponse": {
            "type": "object",
            "properties": {
                "machine": {
                    "$ref": "#/definitions/model.SlotMachine"
                }
            }
        },
        "usecase.PlayRequest": {
            "type": "object",
            "properties": {
                "amount_bet": {
                    "type": "integer"
                },
                "machine_id": {
                    "type": "string"
                },
                "player_id": {
                    "type": "string"
                }
            }
        },
        "usecase.PlayResponse": {
            "type": "object",
            "properties": {
                "player_balance": {
                    "type": "integer"
                },
                "result": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "slot_machine_balance": {
                    "type": "integer"
                },
                "win": {
                    "type": "boolean"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
