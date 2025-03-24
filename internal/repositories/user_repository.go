package repositories

import (
	"database/sql"
	"diploma/internal/models"
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
		if err := rows.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Iin, &user.Role, &user.BiometricDataHash, &user.CreatedAt, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	row := r.db.QueryRow("SELECT * FROM public.user WHERE user_id=$1", id)

	var user models.User
	if err := row.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Iin, &user.Role, &user.BiometricDataHash, &user.CreatedAt, &user.Password); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByIin(iin string) (*models.User, error) {
	row := r.db.QueryRow("SELECT * FROM public.user WHERE iin=$1", iin)

	var user models.User
	if err := row.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Iin, &user.Role, &user.BiometricDataHash, &user.CreatedAt, &user.Password); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	row := r.db.QueryRow("SELECT * FROM public.user WHERE email=$1", email)

	var user models.User
	if err := row.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Iin, &user.Role, &user.BiometricDataHash, &user.CreatedAt, &user.Password); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) DeleteUser(id int) error {
	_, err := r.db.Exec("DELETE FROM public.user WHERE user_id = $1", id)
	return err
}

func (r *UserRepository) CreateUser(user *models.UserRequest) error {
	err := r.db.QueryRow("INSERT INTO public.user (first_name, last_name, email, phone_number, iin, role, password) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING user_id", user.FirstName, user.LastName, user.Email, user.PhoneNumber, user.Iin, user.Role, user.Password).Scan(&user.UserId)
	return err
}

func (r *UserRepository) CreatePatient(patient *models.PatientDetails) error {
	err := r.db.QueryRow("INSERT INTO public.patient (user_id, date_of_birth, gender) VALUES ($1, $2, $3) RETURNING patient_id", patient.UserId, patient.DateOfBirth, patient.Gender).Scan(&patient.PatientId)
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

//func (r *BookRepository) UpdateBook(book *models.Book) error {
//	_, err := r.db.Exec("UPDATE books SET title = $1, author = $2 WHERE id = $3", book.Title, book.Author, book.ID)
//	return err
//}
//
