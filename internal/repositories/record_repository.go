package repositories

import (
	"database/sql"
	"diploma/internal/models"
)

type RecordRepository struct {
	db *sql.DB
}

func NewRecordRepository(db *sql.DB) *RecordRepository {
	return &RecordRepository{db: db}
}

func (r *RecordRepository) GetRecordByUserID(patientId int) (*models.Record, error) {
	row := r.db.QueryRow("SELECT * FROM public.medical_record WHERE patient_id=$1", patientId)

	var record models.Record
	if err := row.Scan(&record.RecordId, &record.PatientId, &record.DoctorId, &record.Diagnosis, &record.TreatmentPlan, &record.TestResult, &record.CreatedAt); err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *RecordRepository) CreateRecord(record *models.Record) error {
	err := r.db.QueryRow("INSERT INTO public.medical_record(patient_id, doctor_id, diagnosis, treatment_plan, test_result) VALUES ($1, $2, $3, $4, $5) RETURNING record_id", record.PatientId, record.DoctorId, record.Diagnosis, record.TreatmentPlan, record.TestResult).Scan(&record.RecordId)
	return err
}
