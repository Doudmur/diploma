package handlers

import (
	"diploma/internal/models"
	"diploma/internal/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type AppointmentHandler struct {
	AppointmentRepo *repositories.AppointmentRepository
	UserRepo        *repositories.UserRepository
	PatientRepo     *repositories.PatientRepository
}

func NewAppointmentHandler(appointmentRepo *repositories.AppointmentRepository, userRepo *repositories.UserRepository, patientRepo *repositories.PatientRepository) *AppointmentHandler {
	return &AppointmentHandler{
		AppointmentRepo: appointmentRepo,
		UserRepo:        userRepo,
		PatientRepo:     patientRepo,
	}
}

// CreateAppointment godoc
// @Summary      Create a new appointment
// @Description  Create a new appointment
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Param        request  body  models.AppointmentRequest  true  "Appointment"
// @Param        Authorization header string true "Bearer"
// @Success      201  {object}  models.Appointment
// @Failure      400  {object}  map[string]string
// @Router       /appointments [post]
func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	var request models.AppointmentRequest
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

	// Validate appointment date
	if request.Date.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Appointment date cannot be in the past"})
		return
	}

	appointment := models.Appointment{
		PatientID: patient.PatientId,
		DoctorID:  doctor.DoctorId,
		Date:      request.Date,
	}

	// Create appointment
	if err := h.AppointmentRepo.CreateAppointment(&appointment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, appointment)
}

// DeleteAppointment godoc
// @Summary      Delete an appointment
// @Description  Delete an appointment if it is before the scheduled date
// @Tags         appointments
// @Produce      json
// @Param        id  path  int  true  "Appointment ID"
// @Param        Authorization header string true "Bearer"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /appointments/{id} [delete]
func (h *AppointmentHandler) DeleteAppointment(c *gin.Context) {
	userID := c.GetUint("user_id")
	doctor, err := h.UserRepo.GetDoctorByUserId(strconv.Itoa(int(userID)))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Doctor not found"})
		return
	}

	appointmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	appointment, err := h.AppointmentRepo.GetAppointmentByID(appointmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	// Ensure the doctor owns the appointment
	if appointment.DoctorID != doctor.DoctorId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to delete this appointment"})
		return
	}

	// Check if the appointment is in the future
	if !appointment.Date.After(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete past or current appointments"})
		return
	}

	if err := h.AppointmentRepo.DeleteAppointment(appointmentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment deleted successfully"})
}

// GetAppointments godoc
// @Summary      Get appointments
// @Description  Fetch appointments for the authenticated doctor or patient
// @Tags         appointments
// @Produce      json
// @Param        Authorization header string true "Bearer"
// @Success      200  {array}  models.Appointment
// @Failure      401  {object}  map[string]string
// @Router       /appointments [get]
func (h *AppointmentHandler) GetAppointments(c *gin.Context) {
	userID := c.GetUint("user_id")
	role, _ := c.Get("role")

	switch role {
	case "doctor":
		doctor, err := h.UserRepo.GetDoctorByUserId(strconv.Itoa(int(userID)))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Doctor not found"})
			return
		}
		appointments, err := h.AppointmentRepo.GetAppointmentsByDoctorID(doctor.DoctorId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, appointments)

	case "patient":
		patient, err := h.PatientRepo.GetPatientByUserID(int(userID))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Patient not found"})
			return
		}
		appointments, err := h.AppointmentRepo.GetAppointmentsByPatientID(patient.PatientId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, appointments)

	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role"})
	}
}
