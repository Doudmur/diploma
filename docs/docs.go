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
        "/auth/login": {
            "post": {
                "description": "Authenticate user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Authentication",
                "parameters": [
                    {
                        "description": "Login request object",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.LoginRequest"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Create a new user with the provided details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Create a new user",
                "parameters": [
                    {
                        "description": "User request object",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.UserRequest"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/notifications": {
            "post": {
                "description": "Create a new notification with the provided details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notifications"
                ],
                "summary": "Create a new notification",
                "parameters": [
                    {
                        "description": "Notification object",
                        "name": "notification",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Notification"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Bearer",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Notification"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/notifications/{id}": {
            "get": {
                "description": "Fetch a notification by its UserID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notifications"
                ],
                "summary": "Get a notification by UserID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Bearer",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Notification"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a notification by its ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notifications"
                ],
                "summary": "Delete a notification",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Notification ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Bearer",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/patients": {
            "get": {
                "description": "Fetch a list of all patients",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "patients"
                ],
                "summary": "Get all patients",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Patient"
                            }
                        }
                    }
                }
            }
        },
        "/patients/{id}": {
            "get": {
                "description": "Fetch a patient by its ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "patients"
                ],
                "summary": "Get a patient by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Patient ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Patient"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/records": {
            "post": {
                "description": "Create a new record with the provided details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Medical Records"
                ],
                "summary": "Create a new record",
                "parameters": [
                    {
                        "description": "Record object",
                        "name": "record",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Record"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Bearer",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Record"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/records/{id}": {
            "get": {
                "description": "Fetch a record by its UserID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Medical Records"
                ],
                "summary": "Get a record by UserID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Bearer",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Record"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "description": "Fetch a list of all users",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get all users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.User"
                            }
                        }
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "description": "Fetch a user by its ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get a user by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a user by its ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Delete a user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "404": {
                        "description": "Not Found",
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
        "models.DoctorDetails": {
            "type": "object",
            "properties": {
                "doctor_id": {
                    "type": "string"
                },
                "specialization": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "models.LoginRequest": {
            "type": "object",
            "properties": {
                "iin": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "models.Notification": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "notification_id": {
                    "type": "integer"
                },
                "sent_at": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "models.Patient": {
            "type": "object",
            "properties": {
                "date_of_birth": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "patient_id": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "models.PatientDetails": {
            "type": "object",
            "properties": {
                "date_of_birth": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "patient_id": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "models.Record": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "diagnosis": {
                    "type": "string"
                },
                "doctor_id": {
                    "type": "integer"
                },
                "patient_id": {
                    "type": "integer"
                },
                "record_id": {
                    "type": "integer"
                },
                "test_result": {
                    "type": "string"
                },
                "treatment_plan": {
                    "type": "string"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "biometric_data_hash": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "iin": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone_number": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "models.UserRequest": {
            "type": "object",
            "properties": {
                "doctor_details": {
                    "$ref": "#/definitions/models.DoctorDetails"
                },
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "iin": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "patient_details": {
                    "$ref": "#/definitions/models.PatientDetails"
                },
                "phone_number": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                },
                "user_id": {
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
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "DiplomaAPI",
	Description:      "This is a REST API for managing medical organization app.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
