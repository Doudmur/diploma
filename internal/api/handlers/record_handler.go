package handlers

import (
	"diploma/internal/models"
	"diploma/internal/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RecordHandler struct {
	RecordRepo  *repositories.RecordRepository
	UserRepo    *repositories.UserRepository
	PatientRepo *repositories.PatientRepository
}

type CreateRecordRequest struct {
	Record models.Record `json:"record"`
}

func NewRecordHandler(recordRepo *repositories.RecordRepository, userRepo *repositories.UserRepository, patientRepo *repositories.PatientRepository) *RecordHandler {
	return &RecordHandler{RecordRepo: recordRepo, UserRepo: userRepo, PatientRepo: patientRepo}
}

// GetRecordByIIN godoc
// @Summary      Get a record by IIN
// @Description  Fetch a record by its IIN
// @Tags         medical records
// @Produce      json
// @Param        iin  path  string  true  "IIN"
// @Param 		 Authorization header string true "Bearer"
// @Success      200  {array}  models.RecordWithDetails
// @Failure      404  {object}  map[string]string
// @Router       /records/{iin} [get]
func (h *RecordHandler) GetRecordByIIN(c *gin.Context) {
	userID := c.GetUint("user_id")
	//doctor, err := h.UserRepo.GetDoctorByUserId(strconv.Itoa(int(userID)))
	//if err != nil {
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": "Doctor not found"})
	//	return
	//}
	_, err := h.UserRepo.GetDoctorByUserId(strconv.Itoa(int(userID)))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Doctor not found"})
		return
	}

	iin := c.Param("iin")
	user, err := h.UserRepo.GetUserByIin(iin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	//patient, err := h.PatientRepo.GetPatientByUserID(user.UserId)
	//if err != nil {
	//	c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
	//	return
	//}
	_, err = h.PatientRepo.GetPatientByUserID(user.UserId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	//// Check if doctor has valid access
	//hasAccess, err := h.RecordRepo.HasValidAccess(doctor.DoctorId, patient.PatientId)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	//
	//if !hasAccess {
	//	c.JSON(http.StatusForbidden, gin.H{"error": "No valid access to patient records"})
	//	return
	//}

	records, err := h.RecordRepo.GetRecordsByIIN(iin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Records not found"})
		return
	}
	c.JSON(http.StatusOK, records)
}

// GetRecordByClaim godoc
// @Summary      Get a record by UserID through claim
// @Description  Fetch a record by its UserID
// @Tags         medical records
// @Produce      json
// @Param 		 Authorization header string true "Bearer"
// @Success      200  {array}  models.RecordWithDetails
// @Failure      404  {object}  map[string]string
// @Router       /records [get]
func (h *RecordHandler) GetRecordByClaim(c *gin.Context) {
	userID := c.GetUint("user_id")
	records, err := h.RecordRepo.GetRecordsByPatientID(int(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Records not found"})
		return
	}
	c.JSON(http.StatusOK, records)
}

// CreateRecord godoc
// @Summary      Create a new record
// @Description  Create a new record
// @Tags         medical records
// @Accept       json
// @Produce      json
// @Param        request  body  models.RecordRequest  true  "Record"
// @Param        Authorization header string true "Bearer"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Router       /records [post]
func (h *RecordHandler) CreateRecord(c *gin.Context) {
	var request models.RecordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	// Get user by IIN
	user, err := h.UserRepo.GetUserByIin(request.Iin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get patient by userId
	patient, err := h.PatientRepo.GetPatientByUserID(user.UserId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	// Get doctor by userId
	doctor, err := h.UserRepo.GetDoctorByUserId(strconv.Itoa(int(userID)))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Doctor not found"})
		return
	}

	record := models.Record{
		PatientId:     patient.PatientId,
		DoctorId:      doctor.DoctorId,
		Iin:           request.Iin,
		Diagnosis:     request.Diagnosis,
		TreatmentPlan: request.TreatmentPlan,
		TestResult:    request.TestResult,
	}

	// Create record
	if err := h.RecordRepo.CreateRecord(&record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var accessLog = models.AccessLog{
		DoctorId:   record.DoctorId,
		RecordId:   record.RecordId,
		AccessType: "CreateRecord",
	}

	// Create access log
	if err := h.RecordRepo.CreateAccessLog(&accessLog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := map[string]interface{}{
		"record": record,
	}

	c.JSON(http.StatusCreated, response)
}

// UpdateRecord godoc
// @Summary      Update a record
// @Description  Update an existing medical record with the provided details
// @Tags         medical records
// @Accept       json
// @Produce      json
// @Param        id   path      int                  true  "Record ID"
// @Param        record  body      models.RecordRequest  true  "Updated Record object"
// @Param 		 Authorization header string true "Bearer"
// @Success      200  {object}  models.Record
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /records/{id} [put]
func (h *RecordHandler) UpdateRecord(c *gin.Context) {
	recordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	var request models.RecordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the doctor ID from the context (set by AuthMiddleware)
	userId := c.GetUint("user_id")
	doctor, err := h.UserRepo.GetDoctorByUserId(strconv.Itoa(int(userId)))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Doctor not found"})
		return
	}

	// Verify the doctor has permission
	existingRecord, err := h.RecordRepo.GetRecordByID(recordID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	if existingRecord.DoctorId != doctor.DoctorId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to update this record"})
		return
	}

	updatedRecord := models.Record{
		RecordId:      recordID,
		PatientId:     existingRecord.PatientId,
		DoctorId:      doctor.DoctorId,
		Iin:           request.Iin,
		Diagnosis:     request.Diagnosis,
		TreatmentPlan: request.TreatmentPlan,
		TestResult:    request.TestResult,
		CreatedAt:     existingRecord.CreatedAt,
	}

	// Update the record
	if err := h.RecordRepo.UpdateRecord(&updatedRecord); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log the update
	var accessLog = models.AccessLog{
		DoctorId:   doctor.DoctorId,
		RecordId:   updatedRecord.RecordId,
		AccessType: "UpdateRecord",
	}

	if err := h.RecordRepo.CreateAccessLog(&accessLog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedRecord)
}
