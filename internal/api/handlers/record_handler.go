package handlers

import (
	"diploma/internal/models"
	"diploma/internal/repositories"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RecordHandler struct {
	repo *repositories.RecordRepository
}

func NewRecordHandler(repo *repositories.RecordRepository) *RecordHandler {
	return &RecordHandler{repo: repo}
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
	record, err := h.repo.GetRecordByIIN(iin)
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
	//if !exists {
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": "ID not found in context"})
	//	c.Abort()
	//	return
	//}
	fmt.Print(userID)
	record, err := h.repo.GetRecordByPatientID(int(userID))
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

	if err := h.repo.CreateRecord(&record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, record)
}
