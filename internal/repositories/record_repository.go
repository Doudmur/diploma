package repositories

import (
	"database/sql"
	"diploma/internal/blockchain"
	"diploma/internal/models"
	"encoding/json"
)

type RecordRepository struct {
	db         *sql.DB
	blockchain *blockchain.Blockchain
}

func NewRecordRepository(db *sql.DB) *RecordRepository {
	return &RecordRepository{db: db, blockchain: blockchain.NewBlockchain(db)}
}
func (r *RecordRepository) GetRecordByPatientID(UserId int) (*models.Record, error) {
	row := r.db.QueryRow("SELECT * FROM public.patient WHERE user_id=$1", UserId)

	var patient models.Patient
	if err := row.Scan(&patient.PatientId, &patient.UserId, &patient.DateOfBirth, &patient.Gender); err != nil {
		return nil, err
	}
	row = r.db.QueryRow("SELECT * FROM public.medical_record WHERE patient_id=$1", patient.PatientId)
	var record models.Record
	if err := row.Scan(&record.RecordId, &record.PatientId, &record.DoctorId, &record.Diagnosis, &record.TreatmentPlan, &record.TestResult, &record.CreatedAt); err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *RecordRepository) GetRecordByIIN(iin string) (*models.Record, error) {
	var userId int
	if err := r.db.QueryRow("SELECT user_id FROM public.user WHERE iin=$1", iin).Scan(&userId); err != nil {
		return nil, err
	}

	var patientId int
	if err := r.db.QueryRow("SELECT patient_id FROM public.patient WHERE user_id=$1", userId).Scan(&patientId); err != nil {
		return nil, err
	}

	row := r.db.QueryRow("SELECT * FROM public.medical_record WHERE patient_id=$1", patientId)
	var record models.Record
	if err := row.Scan(&record.RecordId, &record.PatientId, &record.DoctorId, &record.Diagnosis, &record.TreatmentPlan, &record.TestResult, &record.CreatedAt); err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *RecordRepository) GetRecordByID(recordID int) (*models.Record, error) {
	row := r.db.QueryRow("SELECT * FROM public.medical_record WHERE record_id=$1", recordID)
	var record models.Record
	if err := row.Scan(&record.RecordId, &record.PatientId, &record.DoctorId, &record.Diagnosis, &record.TreatmentPlan, &record.TestResult, &record.CreatedAt); err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *RecordRepository) UpdateRecord(record *models.Record) error {
	// Convert updated record to JSON string for blockchain
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return err
	}

	// Update in database
	_, err = r.db.Exec("UPDATE public.medical_record SET diagnosis=$1, treatment_plan=$2, test_result=$3 WHERE record_id=$4",
		record.Diagnosis, record.TreatmentPlan, record.TestResult, record.RecordId)
	if err != nil {
		return err
	}

	// Add to blockchain
	r.blockchain.AddBlock("Update", record.RecordId, record.DoctorId, record.PatientId, string(recordJSON))

	return nil
}

func (r *RecordRepository) CreateRecord(record *models.Record) error {
	// Convert record to JSON string for blockchain
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return err
	}

	// Store in database
	err = r.db.QueryRow("INSERT INTO public.medical_record(patient_id, doctor_id, diagnosis, treatment_plan, test_result) VALUES ($1, $2, $3, $4, $5) RETURNING record_id", record.PatientId, record.DoctorId, record.Diagnosis, record.TreatmentPlan, record.TestResult).Scan(&record.RecordId)
	if err != nil {
		return err
	}

	// Add to blockchain
	r.blockchain.AddBlock("Create", record.RecordId, record.DoctorId, record.PatientId, string(recordJSON))

	return nil
}

func (r *RecordRepository) CreateAccessLog(accessLog *models.AccessLog) error {
	err := r.db.QueryRow("INSERT INTO access_log(doctor_id, record_id, access_type) VALUES ($1, $2, $3) RETURNING log_id", accessLog.DoctorId, accessLog.RecordId, accessLog.AccessType).Scan(&accessLog.LogId)
	return err
}
