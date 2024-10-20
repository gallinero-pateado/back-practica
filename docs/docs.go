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
        "/check-postulaciones": {
            "get": {
                "description": "Verifica cambios en la tabla de postulaciones y envía correos electrónicos si hay cambios",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "postulaciones"
                ],
                "summary": "Verificar cambios en la tabla de postulaciones",
                "responses": {
                    "200": {
                        "description": "Verificación completada",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Error al verificar cambios",
                        "schema": {
                            "$ref": "#/definitions/database.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/complete-profile": {
            "post": {
                "description": "Permite a los usuarios autenticados completar o actualizar su perfil, incluida la foto de perfil",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "profile"
                ],
                "summary": "Completar o actualizar perfil de usuario",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Datos para actualizar el perfil",
                        "name": "profile",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.ProfileUpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Perfil actualizado correctamente",
                        "schema": {
                            "$ref": "#/definitions/auth.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Datos inválidos",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Usuario no autenticado",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Error al actualizar el perfil",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Autentica al usuario utilizando Firebase y devuelve un token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Inicia sesión un usuario",
                "parameters": [
                    {
                        "description": "Datos de inicio de sesión",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Inicio de sesión exitoso",
                        "schema": {
                            "$ref": "#/definitions/auth.LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Datos inválidos",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Credenciales incorrectas",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/password-reset": {
            "post": {
                "description": "Permite a los usuarios recuperar su contraseña mediante un correo de recuperación",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "password"
                ],
                "summary": "Envía un correo de recuperación de contraseña",
                "parameters": [
                    {
                        "description": "Correo del usuario",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.PasswordResetRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Correo de recuperación enviado con éxito",
                        "schema": {
                            "$ref": "#/definitions/auth.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Email requerido",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Error al enviar el correo de recuperación",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/profile-status": {
            "get": {
                "description": "Retorna si el perfil ha sido completado o no",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "profile"
                ],
                "summary": "Obtener estado del perfil",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Estado del perfil",
                        "schema": {
                            "$ref": "#/definitions/auth.ProfileStatusResponse"
                        }
                    },
                    "400": {
                        "description": "Datos inválidos",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Usuario no autenticado",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Error interno del servidor",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "Crea un nuevo usuario en Firebase y lo guarda en la base de datos local",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Registra un nuevo usuario",
                "parameters": [
                    {
                        "description": "Datos del usuario a registrar",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Usuario registrado correctamente",
                        "schema": {
                            "$ref": "#/definitions/auth.RegisterResponse"
                        }
                    },
                    "400": {
                        "description": "Solicitud inválida",
                        "schema": {
                            "$ref": "#/definitions/auth.RegisterResponse"
                        }
                    },
                    "500": {
                        "description": "Error interno del servidor",
                        "schema": {
                            "$ref": "#/definitions/auth.RegisterResponse"
                        }
                    }
                }
            }
        },
        "/register_admin": {
            "post": {
                "description": "Crea un nuevo usuario en Firebase y lo guarda en la base de datos local",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Registra un nuevo usuario",
                "parameters": [
                    {
                        "description": "Datos del usuario a registrar",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.RegisterRequest_admin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Usuario registrado correctamente",
                        "schema": {
                            "$ref": "#/definitions/auth.RegisterResponse_admin"
                        }
                    },
                    "400": {
                        "description": "Solicitud inválida",
                        "schema": {
                            "$ref": "#/definitions/auth.RegisterResponse_admin"
                        }
                    },
                    "500": {
                        "description": "Error interno del servidor",
                        "schema": {
                            "$ref": "#/definitions/auth.RegisterResponse_admin"
                        }
                    }
                }
            }
        },
        "/register_empresa": {
            "post": {
                "description": "Crea un nuevo usuario en Firebase y lo guarda en la base de datos local",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Registra un nuevo usuario",
                "parameters": [
                    {
                        "description": "Datos del usuario a registrar",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.RegisterRequest_empresa"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Usuario registrado correctamente",
                        "schema": {
                            "$ref": "#/definitions/auth.RegisterResponse_empresa"
                        }
                    },
                    "400": {
                        "description": "Solicitud inválida",
                        "schema": {
                            "$ref": "#/definitions/auth.RegisterResponse_empresa"
                        }
                    },
                    "500": {
                        "description": "Error interno del servidor",
                        "schema": {
                            "$ref": "#/definitions/auth.RegisterResponse_empresa"
                        }
                    }
                }
            }
        },
        "/resend-verification": {
            "post": {
                "description": "Reenvía el correo de verificación a un usuario registrado",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "verification"
                ],
                "summary": "Reenviar correo de verificación",
                "parameters": [
                    {
                        "description": "Correo del usuario",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.EmailRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Correo de verificación enviado nuevamente",
                        "schema": {
                            "$ref": "#/definitions/auth.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Email requerido",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Usuario no encontrado",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Error interno del servidor",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sendemail": {
            "post": {
                "description": "Sends an email to the user with the updated status of their application",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "email"
                ],
                "summary": "Send an email notification",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Recipient email address",
                        "name": "to",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Application status",
                        "name": "estadoPostulacion",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Email sent successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/upload-image": {
            "post": {
                "description": "Sube una imagen a Firebase Storage y actualiza el campo de foto de perfil del usuario autenticado",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "upload"
                ],
                "summary": "Subir una imagen de perfil",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Imagen a subir",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "URL de la imagen subida y mensaje de éxito",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Error en la solicitud",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Usuario no autenticado",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Error al subir la imagen",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.EmailRequest": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "auth.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "auth.LoginRequest": {
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
                    "type": "string"
                }
            }
        },
        "auth.LoginResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                },
                "uid": {
                    "type": "string"
                }
            }
        },
        "auth.PasswordResetRequest": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "auth.ProfileStatusResponse": {
            "type": "object",
            "properties": {
                "perfil_completado": {
                    "type": "boolean"
                }
            }
        },
        "auth.ProfileUpdateRequest": {
            "type": "object",
            "properties": {
                "ano_ingreso": {
                    "type": "string"
                },
                "fecha_nacimiento": {
                    "type": "string"
                },
                "id_carrera": {
                    "type": "integer"
                }
            }
        },
        "auth.RegisterRequest": {
            "type": "object",
            "required": [
                "apellidos",
                "email",
                "nombres",
                "password"
            ],
            "properties": {
                "apellidos": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "nombres": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "auth.RegisterRequest_admin": {
            "type": "object",
            "required": [
                "Email_admin",
                "Nombre_admin",
                "password"
            ],
            "properties": {
                "Email_admin": {
                    "type": "string"
                },
                "Nombre_admin": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "auth.RegisterRequest_empresa": {
            "type": "object",
            "required": [
                "Email_empresa",
                "Nombre_empresa",
                "password"
            ],
            "properties": {
                "Email_empresa": {
                    "type": "string"
                },
                "Nombre_empresa": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "auth.RegisterResponse": {
            "type": "object",
            "properties": {
                "firebase_uid": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "auth.RegisterResponse_admin": {
            "type": "object",
            "properties": {
                "firebase_uid": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "auth.RegisterResponse_empresa": {
            "type": "object",
            "properties": {
                "firebase_uid": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "auth.SuccessResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "database.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
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
