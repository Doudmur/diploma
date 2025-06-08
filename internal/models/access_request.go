package models

import (
	"time"
)

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
