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
}

type Patient struct {
	PatientId   int    `json:"patient_id"`
	UserId      int    `json:"user_id"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
}
