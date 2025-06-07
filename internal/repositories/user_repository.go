package repositories

import (
	"database/sql"
	"diploma/internal/models"
	"fmt"
	"strconv"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUsers() ([]models.User, error) {
	rows, err := r.db.Query("SELECT * FROM public.user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Iin, &user.Role, &user.BiometricDataHash, &user.CreatedAt, &user.Password, &user.PasswordChanged); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	fmt.Print(id)
	row := r.db.QueryRow("SELECT * FROM public.user WHERE user_id=$1", id)

	var user models.User
	if err := row.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Iin, &user.Role, &user.BiometricDataHash, &user.CreatedAt, &user.Password, &user.PasswordChanged); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByIin(iin string) (*models.User, error) {

	row := r.db.QueryRow("SELECT * FROM public.user WHERE iin=$1", iin)
	var user models.User
	if err := row.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Iin, &user.Role, &user.BiometricDataHash, &user.CreatedAt, &user.Password, &user.PasswordChanged, &user.Gender); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	row := r.db.QueryRow("SELECT * FROM public.user WHERE email=$1", email)

	var user models.User
	if err := row.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Iin, &user.Role, &user.BiometricDataHash, &user.CreatedAt, &user.Password, &user.PasswordChanged, &user.Gender); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetDoctorByUserId(userId string) (*models.Doctor, error) {
	row := r.db.QueryRow("SELECT * FROM public.doctor WHERE user_id=$1", userId)

	var doctor models.Doctor
	if err := row.Scan(&doctor.DoctorId, &doctor.UserId, &doctor.Specialization); err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (r *UserRepository) DeleteUser(id int) error {
	_, err := r.db.Exec("DELETE FROM public.user WHERE user_id = $1", id)
	return err
}

func (r *UserRepository) CreateUser(user *models.UserRequest) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Create user
	err = tx.QueryRow(
		"INSERT INTO public.user (first_name, last_name, email, phone_number, iin, role, password, password_changed, gender) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING user_id",
		user.FirstName, user.LastName, user.Email, user.PhoneNumber, user.Iin, user.Role, user.Password, false, user.Gender,
	).Scan(&user.UserId)
	if err != nil {
		return err
	}

	// Create doctor or patient based on role
	if user.Role == "patient" {
		if user.PatientDetails == nil {
			return fmt.Errorf("patient details are required")
		}
		user.PatientDetails.UserId = user.UserId
		err = tx.QueryRow(
			"INSERT INTO public.patient (user_id, date_of_birth) VALUES ($1, $2) RETURNING patient_id",
			user.PatientDetails.UserId, user.PatientDetails.DateOfBirth,
		).Scan(&user.PatientDetails.PatientId)
		if err != nil {
			return err
		}
	} else if user.Role == "doctor" {
		if user.DoctorDetails == nil {
			return fmt.Errorf("doctor details are required")
		}
		user.DoctorDetails.UserId = user.UserId
		err = tx.QueryRow(
			"INSERT INTO public.doctor (user_id, specialization) VALUES ($1, $2) RETURNING doctor_id",
			user.DoctorDetails.UserId, user.DoctorDetails.Specialization,
		).Scan(&user.DoctorDetails.DoctorId)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *UserRepository) CreatePatient(patient *models.PatientDetails) error {
	err := r.db.QueryRow("INSERT INTO public.patient (user_id, date_of_birth) VALUES ($1, $2) RETURNING patient_id", patient.UserId, patient.DateOfBirth).Scan(&patient.PatientId)
	return err
}

func (r *UserRepository) CreateDoctor(doctor *models.DoctorDetails) error {
	err := r.db.QueryRow("INSERT INTO public.doctor (user_id, specialization) VALUES ($1, $2) RETURNING doctor_id", doctor.UserId, doctor.Specialization).Scan(&doctor.DoctorId)
	return err
}

func (r *UserRepository) DeletePatient(id int) error {
	_, err := r.db.Exec("DELETE FROM public.patient WHERE user_id = $1", id)
	return err
}

func (r *UserRepository) DeleteDoctor(id int) error {
	_, err := r.db.Exec("DELETE FROM public.doctor WHERE user_id = $1", id)
	return err
}

func (r *UserRepository) UpdatePassword(iin string, password string) error {
	_, err := r.db.Exec("UPDATE public.user SET password = $1, password_changed = true WHERE iin = $2", password, iin)
	return err
}

func (r *UserRepository) StoreOTP(iin string, otp string, expiresAt time.Time) error {
	_, err := r.db.Exec("INSERT INTO public.otp_verification (iin, otp, expires_at) VALUES ($1, $2, $3) ON CONFLICT (iin) DO UPDATE SET otp = $2, expires_at = $3", iin, otp, expiresAt)
	return err
}

func (r *UserRepository) VerifyOTP(iin string, otp string) (bool, error) {
	var storedOTP string
	var expiresAt time.Time
	err := r.db.QueryRow("SELECT otp, expires_at FROM public.otp_verification WHERE iin = $1", iin).Scan(&storedOTP, &expiresAt)
	if err != nil {
		return false, err
	}

	if time.Now().After(expiresAt) {
		return false, nil
	}

	return storedOTP == otp, nil
}

func (r *UserRepository) DeleteOTP(iin string) error {
	_, err := r.db.Exec("DELETE FROM public.otp_verification WHERE iin = $1", iin)
	return err
}

func (r *UserRepository) ResetPassword(iin string, newPassword string) error {
	_, err := r.db.Exec("UPDATE public.user SET password = $1, password_changed = false WHERE iin = $2", newPassword, iin)
	return err
}

func (r *UserRepository) SetPasswordChanged(iin string, changed bool) error {
	_, err := r.db.Exec("UPDATE public.user SET password_changed = $1 WHERE iin = $2", changed, iin)
	return err
}

// CreateOTPVerification creates a new OTP verification record
func (r *UserRepository) CreateOTPVerification(otp *models.OTPVerification) error {
	query := `
		INSERT INTO otp_verifications (iin, otp, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(query, otp.Iin, otp.OTP, otp.ExpiresAt)
	return err
}

// GetOTPVerification retrieves an OTP verification record by IIN
func (r *UserRepository) GetOTPVerification(iin string) (*models.OTPVerification, error) {
	query := `
		SELECT iin, otp, expires_at
		FROM otp_verifications
		WHERE iin = $1
	`
	var otp models.OTPVerification
	err := r.db.QueryRow(query, iin).Scan(&otp.Iin, &otp.OTP, &otp.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

// DeleteOTPVerification deletes an OTP verification record by IIN
func (r *UserRepository) DeleteOTPVerification(iin string) error {
	query := `
		DELETE FROM otp_verifications
		WHERE iin = $1
	`
	_, err := r.db.Exec(query, iin)
	return err
}

// GetDoctorByID retrieves a doctor by their ID
func (r *UserRepository) GetDoctorByID(doctorID int) (*models.Doctor, error) {
	row := r.db.QueryRow("SELECT * FROM public.doctor WHERE doctor_id=$1", doctorID)

	var doctor models.Doctor
	if err := row.Scan(&doctor.DoctorId, &doctor.UserId, &doctor.Specialization); err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (r *UserRepository) GetUserInfoByIIN(iin string) (*models.UserInfoResponse, error) {
	// Get user information
	user, err := r.GetUserByIin(iin)

	if err != nil {
		return nil, err
	}

	response := &models.UserInfoResponse{
		User: *user,
	}
	// Get role-specific information
	switch user.Role {
	case "doctor":
		doctor, err := r.GetDoctorByUserId(strconv.Itoa(user.UserId))
		if err != nil {
			return nil, err
		}
		response.DoctorDetails = doctor
	case "patient":
		patient, err := r.GetPatientByUserID(user.UserId)
		if err != nil {
			return nil, err
		}
		response.PatientDetails = patient
	}

	return response, nil
}

func (r *UserRepository) GetPatientByUserID(userID int) (*models.Patient, error) {
	row := r.db.QueryRow("SELECT * FROM public.patient WHERE user_id=$1", userID)

	var patient models.Patient
	if err := row.Scan(&patient.PatientId, &patient.UserId, &patient.DateOfBirth); err != nil {
		return nil, err
	}
	return &patient, nil
}

//func (r *BookRepository) UpdateBook(book *models.Book) error {
//	_, err := r.db.Exec("UPDATE books SET title = $1, author = $2 WHERE id = $3", book.Title, book.Author, book.ID)
//	return err
//}
//
