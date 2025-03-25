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
// @Description  Create a new record with the provided details
// @Tags         medical records
// @Accept       json
// @Produce      json
// @Param        record  body  models.Record  true  "Record object"
// @Param 		 Authorization header string true "Bearer"
// @Success      201  {object}  models.Record
// @Failure      400  {object}  map[string]string
// @Router       /records [post]
func (h *RecordHandler) CreateRecord(c *gin.Context) {
	var record models.Record
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := c.GetUint("user_id")

	//Get user by iin
	user, err := h.UserRepo.GetUserByIin(record.Iin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
	}

	// Get patient by userId
	patient, err := h.PatientRepo.GetPatientByUserID(user.UserId)

	// Get doctor by userId
	doctor, err := h.UserRepo.GetDoctorByUserId(strconv.Itoa(int(userId)))

	record.PatientId = patient.PatientId
	record.DoctorId = doctor.DoctorId

	fmt.Print(record.PatientId, record.DoctorId)

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

	c.JSON(http.StatusCreated, record)
}
