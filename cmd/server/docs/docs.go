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
        "/api/auth/singin": {
            "post": {
                "description": "SingIn",
                "tags": [
                    "Auth"
                ],
                "summary": "SingIn",
                "parameters": [
                    {
                        "description": "SingIn request body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.SingInReqBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "SingIn success",
                        "schema": {
                            "$ref": "#/definitions/auth.SingInResp200Body"
                        }
                    },
                    "401": {
                        "description": "invalid credentials",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/auth/singup": {
            "post": {
                "description": "SingUp",
                "tags": [
                    "Auth"
                ],
                "summary": "SingUp",
                "parameters": [
                    {
                        "description": "SingUp request body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.SingUpReqBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "SingUp success",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "409": {
                        "description": "user already exists",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/score/delete": {
            "delete": {
                "security": [
                    {
                        "UserTokenAuth": []
                    }
                ],
                "description": "Deletes the entire leaderboard (for admin use).",
                "tags": [
                    "Score"
                ],
                "summary": "DeleteAllScores",
                "responses": {
                    "200": {
                        "description": "DeleteAllScores success",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/score/list": {
            "post": {
                "description": "Get list of scores based on season",
                "tags": [
                    "Score"
                ],
                "summary": "ListScores",
                "parameters": [
                    {
                        "description": "ListScores request body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/score.ListScoresReqBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ListScores success",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/score.Score"
                            }
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/score/listen_list": {
            "get": {
                "description": "Listen to the score list updates",
                "tags": [
                    "Score"
                ],
                "summary": "ListenScores",
                "responses": {
                    "200": {
                        "description": "ListenScores success",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/score/submit": {
            "post": {
                "security": [
                    {
                        "UserTokenAuth": []
                    }
                ],
                "description": "Create or updates if exists the player’s score",
                "tags": [
                    "Score"
                ],
                "summary": "SubmitScore",
                "parameters": [
                    {
                        "description": "SubmitScore request body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/score.SubmitScoreReqBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "SubmitScore success",
                        "schema": {
                            "$ref": "#/definitions/score.SubmitScoreResp200Body"
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/score/top": {
            "get": {
                "description": "Retrieves the top players based on their ranking",
                "tags": [
                    "Score"
                ],
                "summary": "TopScores",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "season",
                        "name": "season",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "TopScores success",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/score.Score"
                            }
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.SingInReqBody": {
            "type": "object",
            "required": [
                "nickName",
                "password"
            ],
            "properties": {
                "nickName": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 3
                },
                "password": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 8
                }
            }
        },
        "auth.SingInResp200Body": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "auth.SingUpReqBody": {
            "type": "object",
            "required": [
                "nickName",
                "password"
            ],
            "properties": {
                "nickName": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 3
                },
                "password": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 8
                }
            }
        },
        "score.Details": {
            "type": "object",
            "properties": {
                "losses": {
                    "type": "integer"
                },
                "playerNickname": {
                    "type": "string"
                },
                "rank": {
                    "type": "integer"
                },
                "rating": {
                    "type": "integer"
                },
                "region": {
                    "type": "string"
                },
                "winLoseRatio": {
                    "type": "string"
                },
                "wins": {
                    "type": "integer"
                }
            }
        },
        "score.ListScoresReqBody": {
            "type": "object",
            "required": [
                "season"
            ],
            "properties": {
                "last_received_id": {
                    "type": "string"
                },
                "season": {
                    "type": "string"
                }
            }
        },
        "score.Score": {
            "type": "object",
            "properties": {
                "scoreDetails": {
                    "$ref": "#/definitions/score.Details"
                },
                "scoreID": {
                    "type": "string"
                },
                "scoredAt": {
                    "type": "string"
                },
                "season": {
                    "type": "string"
                }
            }
        },
        "score.SubmitScoreReqBody": {
            "type": "object",
            "required": [
                "losses",
                "rating",
                "region",
                "wins"
            ],
            "properties": {
                "losses": {
                    "type": "integer",
                    "minimum": 0
                },
                "rating": {
                    "description": "NickName string ` + "`" + `json:\"nickName\" validate:\"required,min=3,max=32,alphanum\"` + "`" + `",
                    "type": "integer",
                    "minimum": 1
                },
                "region": {
                    "type": "string",
                    "enum": [
                        "EU",
                        "NA",
                        "AS",
                        "SA",
                        "AF"
                    ]
                },
                "wins": {
                    "type": "integer",
                    "minimum": 0
                }
            }
        },
        "score.SubmitScoreResp200Body": {
            "type": "object",
            "properties": {
                "score": {
                    "$ref": "#/definitions/score.Score"
                }
            }
        }
    },
    "securityDefinitions": {
        "UserTokenAuth": {
            "type": "apiKey",
            "name": "X-User-Token",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Leaderboard API",
	Description:      "This is a sample server for Leaderboard API.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
