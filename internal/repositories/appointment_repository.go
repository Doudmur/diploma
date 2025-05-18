package repositories

import (
	"database/sql"
	"diploma/internal/models"
	"fmt"
)

type AppointmentRepository struct {
	db *sql.DB
}

func NewAppointmentRepository(db *sql.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func (r *AppointmentRepository) CreateAppointment(appointment *models.Appointment) error {
	fmt.Print(appointment)
	err := r.db.QueryRow(`INSERT INTO public.appointments (doctor_id, patient_id, date)VALUES ($1, $2, $3) RETURNING id`,
		appointment.DoctorID, appointment.PatientID, appointment.Date).Scan(&appointment.ID)
	return err
}

func (r *AppointmentRepository) DeleteAppointment(id int) error {
	_, err := r.db.Exec("DELETE FROM public.appointments WHERE id = $1", id)
	return err
}

func (r *AppointmentRepository) GetAppointmentByID(id int) (*models.Appointment, error) {
	row := r.db.QueryRow("SELECT id, doctor_id, patient_id, date FROM public.appointments WHERE id = $1", id)
	var appointment models.Appointment
	if err := row.Scan(&appointment.ID, &appointment.DoctorID, &appointment.PatientID, &appointment.Date); err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *AppointmentRepository) GetAppointmentsByDoctorID(doctorID int) ([]models.Appointment, error) {
	rows, err := r.db.Query("SELECT id, doctor_id, patient_id, date FROM public.appointments WHERE doctor_id = $1", doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		var appointment models.Appointment
		if err := rows.Scan(&appointment.ID, &appointment.DoctorID, &appointment.PatientID, &appointment.Date); err != nil {
			return nil, err
		}
		appointments = append(appointments, appointment)
	}
	return appointments, nil
}

func (r *AppointmentRepository) GetAppointmentsByPatientID(patientID int) ([]models.Appointment, error) {
	rows, err := r.db.Query("SELECT id, doctor_id, patient_id, date FROM public.appointments WHERE patient_id = $1", patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		var appointment models.Appointment
		if err := rows.Scan(&appointment.ID, &appointment.DoctorID, &appointment.PatientID, &appointment.Date); err != nil {
			return nil, err
		}
		appointments = append(appointments, appointment)
	}
	return appointments, nil
}
