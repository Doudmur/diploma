package repositories

import (
	"database/sql"
	"diploma/internal/models"
)

type AppointmentRepository struct {
	db *sql.DB
}

func NewAppointmentRepository(db *sql.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func (r *AppointmentRepository) CreateAppointment(appointment *models.Appointment) error {
	err := r.db.QueryRow(
		"INSERT INTO public.appointment (doctor_id, patient_id, date) VALUES ($1, $2, $3) RETURNING id",
		appointment.DoctorID, appointment.PatientID, appointment.Date,
	).Scan(&appointment.ID)
	return err
}

func (r *AppointmentRepository) DeleteAppointment(id int) error {
	_, err := r.db.Exec("DELETE FROM public.appointment WHERE id = $1", id)
	return err
}

func (r *AppointmentRepository) GetAppointmentByID(id int) (*models.Appointment, error) {
	query := `
		SELECT a.id, a.doctor_id, a.patient_id, a.date,
			u.first_name, u.last_name, d.specialization,
			pu.iin as patient_iin
		FROM public.appointment a
		JOIN public.doctor d ON a.doctor_id = d.doctor_id
		JOIN public.user u ON d.user_id = u.user_id
		JOIN public.patient p ON a.patient_id = p.patient_id
		JOIN public.user pu ON p.user_id = pu.user_id
		WHERE a.id = $1`

	var appointment models.Appointment
	if err := r.db.QueryRow(query, id).Scan(
		&appointment.ID,
		&appointment.DoctorID,
		&appointment.PatientID,
		&appointment.Date,
		&appointment.FirstName,
		&appointment.LastName,
		&appointment.Specialization,
		&appointment.Iin,
	); err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *AppointmentRepository) GetAppointmentsByDoctorID(doctorID int) ([]models.Appointment, error) {
	query := `
		SELECT a.id, a.doctor_id, a.patient_id, a.date,
			u.first_name, u.last_name, d.specialization,
			pu.iin as patient_iin
		FROM public.appointment a
		JOIN public.doctor d ON a.doctor_id = d.doctor_id
		JOIN public.user u ON d.user_id = u.user_id
		JOIN public.patient p ON a.patient_id = p.patient_id
		JOIN public.user pu ON p.user_id = pu.user_id
		WHERE a.doctor_id = $1
		ORDER BY a.date DESC`

	rows, err := r.db.Query(query, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		var appt models.Appointment
		if err := rows.Scan(
			&appt.ID,
			&appt.DoctorID,
			&appt.PatientID,
			&appt.Date,
			&appt.FirstName,
			&appt.LastName,
			&appt.Specialization,
			&appt.Iin,
		); err != nil {
			return nil, err
		}
		print(appt.Iin)
		appointments = append(appointments, appt)
	}
	return appointments, nil
}

func (r *AppointmentRepository) GetAppointmentsByPatientID(patientID int) ([]models.Appointment, error) {
	query := `
		SELECT a.id, a.doctor_id, a.patient_id, a.date,
			u.first_name, u.last_name, d.specialization,
			pu.iin as patient_iin
		FROM public.appointment a
		JOIN public.doctor d ON a.doctor_id = d.doctor_id
		JOIN public.user u ON d.user_id = u.user_id
		JOIN public.patient p ON a.patient_id = p.patient_id
		JOIN public.user pu ON p.user_id = pu.user_id
		WHERE a.patient_id = $1
		ORDER BY a.date DESC`

	rows, err := r.db.Query(query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		var appt models.Appointment
		if err := rows.Scan(
			&appt.ID,
			&appt.DoctorID,
			&appt.PatientID,
			&appt.Date,
			&appt.FirstName,
			&appt.LastName,
			&appt.Specialization,
			&appt.Iin,
		); err != nil {
			return nil, err
		}
		appointments = append(appointments, appt)
	}
	return appointments, nil
}
