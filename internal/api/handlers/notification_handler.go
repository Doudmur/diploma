package handlers

import (
	"diploma/internal/models"
	"diploma/internal/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type NotificationHandler struct {
	repo *repositories.NotificationRepository
}

func NewNotificationHandler(repo *repositories.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{repo: repo}
}

// GetNotificationByUserID godoc
// @Summary      Get a notification by UserID
// @Description  Fetch a notification by its UserID
// @Tags         notifications
// @Produce      json
// @Param        id  path  int  true  "User ID"
// @Param 		 Authorization header string true "Bearer"
// @Success      200  {object}  models.Notification
// @Failure      404  {object}  map[string]string
// @Router       /notifications/{id} [get]
func (h *NotificationHandler) GetNotificationByUserID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, err := h.repo.GetNotificationByUserID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteNotification godoc
// @Summary      Delete a notification
// @Description  Delete a notification by its ID
// @Tags         notifications
// @Produce      json
// @Param        id  path  int  true  "Notification ID"
// @Param 		 Authorization header string true "Bearer"
// @Success      204
// @Failure      404  {object}  map[string]string
// @Router       /notifications/{id} [delete]
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.repo.DeleteNotification(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// CreateNotification godoc
// @Summary      Create a new notification
// @Description  Create a new notification with the provided details
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Param        notification  body  models.Notification  true  "Notification object"
// @Param 		 Authorization header string true "Bearer"
// @Success      201  {object}  models.Notification
// @Failure      400  {object}  map[string]string
// @Router       /notifications [post]
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var notification models.Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateNotification(&notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, notification)
}
