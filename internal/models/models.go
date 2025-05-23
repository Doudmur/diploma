package models

import "time"

type User struct {
	UserId            int    `json:"user_id"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email"`
	PhoneNumber       string `json:"phone_number"`
	Iin               string `json:"iin"`
	Role              string `json:"role"`
	BiometricDataHash string `json:"biometric_data_hash"`
	CreatedAt         string `json:"created_at"`
	Password          string `json:"password"`
}

type LoginRequest struct {
	Iin      string `json:"iin"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	Iin            string `json:"iin"`
	OldPassword    string `json:"old_password"`
	NewPassword    string `json:"new_password"`
	VerifyPassword string `json:"verify_password"`
}

type UserRequest struct {
	UserId         int             `json:"user_id"`
	FirstName      string          `json:"first_name"`
	LastName       string          `json:"last_name"`
	Email          string          `json:"email"`
	PhoneNumber    string          `json:"phone_number"`
	Iin            string          `json:"iin"`
	Role           string          `json:"role"`
	Password       string          `json:"password"`
	PatientDetails *PatientDetails `json:"patient_details,omitempty"`
	DoctorDetails  *DoctorDetails  `json:"doctor_details,omitempty"`
}

type PatientDetails struct {
	PatientId   string `json:"patient_id"`
	UserId      int    `json:"user_id"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
}

type DoctorDetails struct {
	DoctorId       string `json:"doctor_id"`
	UserId         int    `json:"user_id"`
	Specialization string `json:"specialization"`
}

type UserResponse struct {
	UserID    int `json:"user_id"`
	PatientID int `json:"patient_id"`
	DoctorID  int `json:"doctor_id"`
}

type Patient struct {
	PatientId   int    `json:"patient_id"`
	UserId      int    `json:"user_id"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
}

type Doctor struct {
	DoctorId       int    `json:"doctor_id"`
	UserId         int    `json:"user_id"`
	Specialization string `json:"specialization"`
}

type Notification struct {
	NotificationId int    `json:"notification_id"`
	UserId         int    `json:"user_id"`
	Message        string `json:"message"`
	Type           string `json:"type"`
	SentAt         string `json:"sent_at"`
}

type Record struct {
	RecordId      int    `json:"record_id"`
	PatientId     int    `json:"patient_id"`
	Iin           string `json:"iin"`
	DoctorId      int    `json:"doctor_id"`
	Diagnosis     string `json:"diagnosis"`
	TreatmentPlan string `json:"treatment_plan"`
	TestResult    string `json:"test_result"`
	CreatedAt     string `json:"created_at"`
}

type AccessLog struct {
	LogId      int    `json:"log_id"`
	DoctorId   int    `json:"doctor_id"`
	RecordId   int    `json:"record_id"`
	AccessType string `json:"access_type"`
	AccessDate string `json:"access_date"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message" example:"Bad request"`
}

type Blockchain struct {
	Chain []Block `json:"chain"`
}

type Block struct {
	Index        int         `json:"index"`
	Timestamp    time.Time   `json:"timestamp"`
	Transaction  Transaction `json:"transaction"`
	PreviousHash string      `json:"previous_hash"`
	Hash         string      `json:"hash"`
}

type Transaction struct {
	Action    string    `json:"action"`     // "Create" or "Update"
	RecordID  int       `json:"record_id"`  // ID of the medical record
	DoctorID  int       `json:"doctor_id"`  // Doctor who performed the action
	PatientID int       `json:"patient_id"` // Patient associated with the record
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details"` // JSON string of record data
}

// Appointment represents an appointment between a doctor and a patient
type Appointment struct {
	ID        int       `json:"id"`
	DoctorID  int       `json:"doctor_id"`
	PatientID int       `json:"patient_id"`
	Iin       string    `json:"iin"`
	Date      time.Time `json:"date"`
}
