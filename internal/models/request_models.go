package models

import "time"

// RecordRequest represents the request for creating/updating a medical record
type RecordRequest struct {
	Iin           string `json:"iin" example:"123456789012"`
	Diagnosis     string `json:"diagnosis" example:"Common cold"`
	TreatmentPlan string `json:"treatment_plan" example:"Rest and medication"`
	TestResult    string `json:"test_result" example:"Blood test results"`
}

// AppointmentRequest represents the request for creating an appointment
type AppointmentRequest struct {
	Iin  string    `json:"iin" example:"123456789012"`
	Date time.Time `json:"date" example:"2024-03-14T12:00:00Z"`
}

// AccessRequestCreate represents the request for creating an access request
type AccessRequestCreate struct {
	Iin string `json:"iin" example:"123456789012"`
}

// AccessRequestStatusUpdate represents the request for updating access request status
type AccessRequestStatusUpdate struct {
	Status string `json:"status" example:"granted" enums:"granted,rejected"`
}

// VerifyOTPRequest represents the request for verifying OTP
type VerifyOTPRequest struct {
	Iin string `json:"iin" example:"123456789012"`
	OTP string `json:"otp" example:"123456"`
}
