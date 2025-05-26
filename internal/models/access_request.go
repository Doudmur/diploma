package models

import (
	"database/sql"
	"time"
)

// AccessRequest represents an access request in the database
type AccessRequest struct {
	ID              int          `json:"id" db:"id"`
	DoctorID        int          `json:"doctor_id" db:"doctor_id"`
	PatientID       int          `json:"patient_id" db:"patient_id"`
	Status          string       `json:"status" db:"status"`
	CreatedAt       time.Time    `json:"created_at" db:"created_at"`
	ExpiresAt       time.Time    `json:"expires_at" db:"expires_at"`
	AccessGrantedAt sql.NullTime `json:"access_granted_at" db:"access_granted_at"`
	AccessExpiresAt sql.NullTime `json:"access_expires_at" db:"access_expires_at"`
}

// AccessRequestResponse represents the response for access request endpoints
type AccessRequestResponse struct {
	ID              int          `json:"id" example:"1"`
	DoctorID        int          `json:"doctor_id" example:"1"`
	DoctorName      string       `json:"doctor_name" example:"John Doe"`
	PatientID       int          `json:"patient_id" example:"1"`
	Status          string       `json:"status" example:"pending" enums:"pending,granted,rejected"`
	CreatedAt       time.Time    `json:"created_at" example:"2024-03-14T12:00:00Z"`
	ExpiresAt       time.Time    `json:"expires_at" example:"2024-03-14T13:00:00Z"`
	AccessGrantedAt sql.NullTime `json:"access_granted_at,omitempty" example:"2024-03-14T12:30:00Z"`
	AccessExpiresAt sql.NullTime `json:"access_expires_at,omitempty" example:"2024-03-14T13:30:00Z"`
}

// AccessRequestResponseSwagger represents the Swagger documentation for access request responses
type AccessRequestResponseSwagger struct {
	ID              int       `json:"id" example:"1"`
	DoctorID        int       `json:"doctor_id" example:"1"`
	DoctorName      string    `json:"doctor_name" example:"John Doe"`
	PatientID       int       `json:"patient_id" example:"1"`
	Status          string    `json:"status" example:"pending" enums:"pending,granted,rejected"`
	CreatedAt       time.Time `json:"created_at" example:"2024-03-14T12:00:00Z"`
	ExpiresAt       time.Time `json:"expires_at" example:"2024-03-14T13:00:00Z"`
	AccessGrantedAt time.Time `json:"access_granted_at,omitempty" example:"2024-03-14T12:30:00Z"`
	AccessExpiresAt time.Time `json:"access_expires_at,omitempty" example:"2024-03-14T13:30:00Z"`
}
