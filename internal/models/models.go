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
	PasswordChanged   bool   `json:"password_changed"`
	Gender            string `json:"gender"`
	Photo             []byte `json:"photo"`
}

type DetectResponse struct {
	HasFace   bool   `json:"has_face"`
	FaceCount int    `json:"face_count"`
	Message   string `json:"message"`
}

type SimilarityResponse struct {
	Distance float64 `json:"distance"`
	Match    bool    `json:"match"`
	Message  string  `json:"message"`
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
	Gender         string          `json:"gender"`
	PatientDetails *PatientDetails `json:"patient_details,omitempty"`
	DoctorDetails  *DoctorDetails  `json:"doctor_details,omitempty"`
}

type PatientDetails struct {
	PatientId   string `json:"patient_id"`
	UserId      int    `json:"user_id"`
	DateOfBirth string `json:"date_of_birth"`
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
}

type Doctor struct {
	DoctorId       int    `json:"doctor_id"`
	UserId         int    `json:"user_id"`
	Specialization string `json:"specialization"`
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
	ID             int       `json:"id"`
	DoctorID       int       `json:"doctor_id"`
	PatientID      int       `json:"patient_id"`
	Date           time.Time `json:"date"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Specialization string    `json:"specialization"`
	Iin            string    `json:"iin"`
}

type ForgotPasswordRequest struct {
	Iin string `json:"iin"`
}

type ResetPasswordRequest struct {
	Iin         string `json:"iin"`
	OTP         string `json:"otp"`
	NewPassword string `json:"new_password"`
}

type OTPVerification struct {
	Iin       string    `json:"iin"`
	OTP       string    `json:"otp"`
	ExpiresAt time.Time `json:"expires_at"`
}

// UserInfoResponse represents detailed user information including role-specific details
type UserInfoResponse struct {
	User           User     `json:"user"`                      // Basic user information
	DoctorDetails  *Doctor  `json:"doctor_details,omitempty"`  // Doctor-specific details if user is a doctor
	PatientDetails *Patient `json:"patient_details,omitempty"` // Patient-specific details if user is a patient
}

type RecordWithDetails struct {
	Record
	DoctorFullName   string `json:"doctor_full_name"`
	DoctorSpeciality string `json:"doctor_specialization"`
	PatientFullName  string `json:"patient_full_name"`
	PatientIIN       string `json:"patient_iin"`
}

type AccessRequest struct {
	ID              int       `json:"id"`
	DoctorID        int       `json:"doctor_id"`
	PatientID       int       `json:"patient_id"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	ExpiresAt       time.Time `json:"expires_at"`
	AccessGrantedAt time.Time `json:"access_granted_at,omitempty"`
	AccessExpiresAt time.Time `json:"access_expires_at,omitempty"`
}

type CreateAccessRequestRequest struct {
	PatientIIN string `json:"patient_iin" binding:"required"`
}

type AccessRequestResponse struct {
	RequestID int       `json:"request_id"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
}
