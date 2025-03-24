package handlers

import (
	"diploma/internal/models"
	"diploma/internal/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RecordHandler struct {
	repo *repositories.RecordRepository
}

func NewRecordHandler(repo *repositories.RecordRepository) *RecordHandler {
	return &RecordHandler{repo: repo}
}

// GetRecordByUserID godoc
// @Summary      Get a record by UserID
// @Description  Fetch a record by its UserID
// @Tags         Medical Records
// @Produce      json
// @Param        id  path  int  true  "User ID"
// @Param 		 Authorization header string true "Bearer"
// @Success      200  {object}  models.Record
// @Failure      404  {object}  map[string]string
// @Router       /records/{id} [get]
func (h *RecordHandler) GetRecordByUserID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	record, err := h.repo.GetRecordByUserID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, record)
}

// CreateRecord godoc
// @Summary      Create a new record
// @Description  Create a new record with the provided details
// @Tags         Medical Records
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
