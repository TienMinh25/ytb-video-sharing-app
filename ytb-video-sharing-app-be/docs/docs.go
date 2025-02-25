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
        "/accounts/check-token": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "check token every user access web page",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "accounts"
                ],
                "summary": "check access token",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.CheckTokenResponseDocs"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/accounts/login": {
            "post": {
                "description": "Authenticate user and return access token \u0026 refresh token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "accounts"
                ],
                "summary": "Login account",
                "parameters": [
                    {
                        "description": "Login payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.LoginResponseDocs"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ResponseError"
                        }
                    }
                }
            }
        },
        "/accounts/logout/{accountID}": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Logout user by deleting refresh token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "accounts"
                ],
                "summary": "Logout account",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Account ID",
                        "name": "accountID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Refresh Token",
                        "name": "X-Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.LogoutResponseDocs"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ResponseError"
                        }
                    }
                }
            }
        },
        "/accounts/refresh-token": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Refresh access token using a valid refresh token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "accounts"
                ],
                "summary": "Refresh Token",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Account ID",
                        "name": "accountID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Refresh Token",
                        "name": "X-Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.RefreshTokenResponseDocs"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/accounts/register": {
            "post": {
                "description": "create new account based info request",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "accounts"
                ],
                "summary": "Register new account",
                "parameters": [
                    {
                        "description": "Thông tin đăng ký",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.CreateAccountRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.CreateAccountResponseDocs"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ResponseError"
                        }
                    }
                }
            }
        },
        "/videos": {
            "get": {
                "description": "Get list videos",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "videos"
                ],
                "summary": "Get list videos",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Limit number of records returned",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "page",
                        "name": "page",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.ListVideosResponseDocs"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ResponseError"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Create new video and return itself.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "videos"
                ],
                "summary": "Share new video",
                "parameters": [
                    {
                        "type": "string",
                        "description": "WebSocket connection ID",
                        "name": "conn_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "description": "Share video payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.ShareVideoRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dto.ShareVideoResponseDocs"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ResponseError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.CheckTokenResponse": {
            "type": "object"
        },
        "dto.CheckTokenResponseDocs": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/dto.CheckTokenResponse"
                },
                "metadata": {
                    "$ref": "#/definitions/dto.Metadata"
                }
            }
        },
        "dto.CreateAccountRequest": {
            "type": "object",
            "required": [
                "email",
                "fullname",
                "password"
            ],
            "properties": {
                "avatar_url": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "fullname": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 6
                }
            }
        },
        "dto.CreateAccountResponseDocs": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/dto.CreateAccountResponseWithOTP"
                },
                "metadata": {
                    "$ref": "#/definitions/dto.Metadata"
                }
            }
        },
        "dto.CreateAccountResponseWithOTP": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "avatar_url": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "fullname": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "otp": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                }
            }
        },
        "dto.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "dto.ListVideosResponseDocs": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.VideoResponse"
                    }
                },
                "metadata": {
                    "$ref": "#/definitions/dto.MetadataWithPagination"
                }
            }
        },
        "dto.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 6
                }
            }
        },
        "dto.LoginResponseDocs": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/dto.LoginResponseWithOTP"
                },
                "metadata": {
                    "$ref": "#/definitions/dto.Metadata"
                }
            }
        },
        "dto.LoginResponseWithOTP": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "avatar_url": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "fullname": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "otp": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                }
            }
        },
        "dto.LogoutResponse": {
            "type": "object"
        },
        "dto.LogoutResponseDocs": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/dto.LogoutResponse"
                },
                "metadata": {
                    "$ref": "#/definitions/dto.Metadata"
                }
            }
        },
        "dto.Metadata": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                }
            }
        },
        "dto.MetadataWithPagination": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "pagination": {
                    "$ref": "#/definitions/dto.Pagination"
                }
            }
        },
        "dto.Pagination": {
            "type": "object",
            "properties": {
                "is_next": {
                    "type": "boolean"
                },
                "is_previous": {
                    "type": "boolean"
                },
                "limit": {
                    "type": "integer"
                },
                "page": {
                    "type": "integer"
                },
                "total_items": {
                    "type": "integer"
                },
                "total_pages": {
                    "type": "integer"
                }
            }
        },
        "dto.RefreshTokenResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                }
            }
        },
        "dto.RefreshTokenResponseDocs": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/dto.RefreshTokenResponse"
                },
                "metadata": {
                    "$ref": "#/definitions/dto.Metadata"
                }
            }
        },
        "dto.ResponseError": {
            "type": "object",
            "properties": {
                "error": {},
                "metadata": {
                    "$ref": "#/definitions/dto.Metadata"
                }
            }
        },
        "dto.ShareVideoRequest": {
            "type": "object",
            "required": [
                "thumbnail",
                "title",
                "video_url"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "downvote": {
                    "type": "integer"
                },
                "thumbnail": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "upvote": {
                    "type": "integer"
                },
                "video_url": {
                    "type": "string"
                }
            }
        },
        "dto.ShareVideoResponse": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "downvote": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "thumbnail": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "upvote": {
                    "type": "integer"
                },
                "video_url": {
                    "type": "string"
                }
            }
        },
        "dto.ShareVideoResponseDocs": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/dto.ShareVideoResponse"
                },
                "metadata": {
                    "$ref": "#/definitions/dto.Metadata"
                }
            }
        },
        "dto.VideoResponse": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "downvote": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "shared_by": {
                    "type": "string"
                },
                "thumbnail": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "upvote": {
                    "type": "integer"
                },
                "video_url": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3000",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "YouTube Video Sharing API",
	Description:      "API cho ứng dụng chia sẻ video YouTube",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
