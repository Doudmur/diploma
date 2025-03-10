package repositories

import (
	"database/sql"
	"diploma/internal/models"
)

type PatientRepository struct {
	db *sql.DB
}

func NewPatientRepository(db *sql.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

func (r *PatientRepository) GetPatients() ([]models.Patient, error) {
	rows, err := r.db.Query("SELECT * FROM public.patient")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []models.Patient
	for rows.Next() {
		var patient models.Patient
		if err := rows.Scan(&patient.PatientId, &patient.UserId, &patient.DateOfBirth, &patient.Gender); err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}
	return patients, nil
}

func (r *PatientRepository) GetPatientByID(id int) (*models.Patient, error) {
	row := r.db.QueryRow("SELECT * FROM public.patient WHERE patient_id=$1", id)

	var patient models.Patient
	if err := row.Scan(&patient.PatientId, &patient.UserId, &patient.DateOfBirth, &patient.Gender); err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *PatientRepository) GetPatientByUserID(id int) (*models.Patient, error) {
	row := r.db.QueryRow("SELECT * FROM public.patient WHERE user_id=$1", id)

	var patient models.Patient
	if err := row.Scan(&patient.PatientId, &patient.UserId, &patient.DateOfBirth, &patient.Gender); err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *PatientRepository) DeletePatient(id int) error {
	_, err := r.db.Exec("DELETE FROM public.patient WHERE patient_id = $1", id)
	return err
}

//func (r *BookRepository) CreateBook(book *models.Book) error {
//	err := r.db.QueryRow("INSERT INTO books (title, author) VALUES ($1, $2) RETURNING id", book.Title, book.Author).Scan(&book.ID)
//	return err
//}
//
//func (r *BookRepository) UpdateBook(book *models.Book) error {
//	_, err := r.db.Exec("UPDATE books SET title = $1, author = $2 WHERE id = $3", book.Title, book.Author, book.ID)
//	return err
//}
//
