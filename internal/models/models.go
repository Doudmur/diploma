package models

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
	Iin string `json:"iin"`

	Password string `json:"password"`
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

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message" example:"Bad request"`
}
