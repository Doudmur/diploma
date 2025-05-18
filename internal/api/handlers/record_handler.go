package handlers

import (
	"diploma/internal/models"
	"diploma/internal/repositories"
	"fmt"
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
// @Success      200  {object}  models.Record
// @Failure      404  {object}  map[string]string
// @Router       /records/{iin} [get]
func (h *RecordHandler) GetRecordByIIN(c *gin.Context) {
	iin := c.Param("iin")
	record, err := h.RecordRepo.GetRecordByIIN(iin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, record)
}

// GetRecordByClaim godoc
// @Summary      Get a record by UserID through claim
// @Description  Fetch a record by its UserID
// @Tags         medical records
// @Produce      json
// @Param 		 Authorization header string true "Bearer"
// @Success      200  {object}  models.Record
// @Failure      404  {object}  map[string]string
// @Router       /records [get]
func (h *RecordHandler) GetRecordByClaim(c *gin.Context) {
	userID := c.GetUint("user_id")
	fmt.Print(userID)
	record, err := h.RecordRepo.GetRecordByPatientID(int(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, record)
}

// CreateRecord godoc
// @Summary      Create a new record
// @Description  Create a new record
// @Tags         medical records
// @Accept       json
// @Produce      json
// @Param        request  body  CreateRecordRequest  true  "Record"
// @Param        Authorization header string true "Bearer"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Router       /records [post]
func (h *RecordHandler) CreateRecord(c *gin.Context) {
	var request CreateRecordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	record := request.Record

	userID := c.GetUint("user_id")

	// Get user by IIN
	user, err := h.UserRepo.GetUserByIin(record.Iin)
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

	record.PatientId = patient.PatientId
	record.DoctorId = doctor.DoctorId

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
// @Param        record  body      models.Record  true  "Updated Record object"
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

	var updatedRecord models.Record
	if err := c.ShouldBindJSON(&updatedRecord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the RecordID from the path parameter
	updatedRecord.RecordId = recordID

	// Get the doctor ID from the context (set by AuthMiddleware)
	userId := c.GetUint("user_id")
	doctor, err := h.UserRepo.GetDoctorByUserId(strconv.Itoa(int(userId)))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Doctor not found"})
		return
	}

	// Verify the doctor has permission (simplified check; in practice, you might check ownership or admin rights)
	existingRecord, err := h.RecordRepo.GetRecordByID(recordID) // Assume this method exists or add it
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	if existingRecord.DoctorId != doctor.DoctorId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to update this record"})
		return
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
