package handlers

import (
	"diploma/internal/auth"
	"diploma/internal/models"
	"diploma/internal/repositories"
	"diploma/internal/scripts"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserHandler struct {
	repo *repositories.UserRepository
}

func NewUserHandler(repo *repositories.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// GetUsers godoc
// @Summary      Get all users
// @Description  Fetch a list of all users
// @Tags         users
// @Produce      json
// @Success      200  {array}  models.User
// @Router       /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.repo.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary      Get a user by ID
// @Description  Fetch a user by its ID
// @Tags         users
// @Produce      json
// @Param        id  path  int  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      404  {object}  map[string]string
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, err := h.repo.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user with the provided details
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body  models.UserRequest  true  "User request object"
// @Success      201  {object}  models.UserRequest
// @Failure      400  {object}  map[string]string
// @Router       /auth/register [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var userRequest models.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid request body. Error": err.Error()})
		return
	}

	// Validate role
	if userRequest.Role != "doctor" && userRequest.Role != "patient" {

		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Invalid role"})
		return
	}

	// Generate OTP
	otp := scripts.GenerateOTP()

	// Send OTP to email
	subject := "Generated password to login in MedicineApp"
	body := "Your password generated to first login: " + otp
	err := scripts.SendMail(userRequest.Email, subject, body)
	if err != nil {
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(otp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	userRequest.Password = hashedPassword

	// Create User
	err = h.repo.CreateUser(&userRequest)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not created"})
		return
	}

	// Create doctor or patient based on role
	if userRequest.Role == "patient" {
		if userRequest.PatientDetails == nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Patient details are required"})
			return
		}
		userRequest.PatientDetails.UserId = userRequest.UserId
		// Create patient
		err = h.repo.CreatePatient(userRequest.PatientDetails)
		if err != nil {
			fmt.Printf("Error creating patient: %v. User ID: %d", err, userRequest.UserId)
			c.JSON(http.StatusNotFound, gin.H{"error": "Patient not created"})
			return
		}
	} else if userRequest.Role == "doctor" {
		if userRequest.DoctorDetails == nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Doctor details are required"})
			return
		}
		userRequest.DoctorDetails.UserId = userRequest.UserId
		// Create doctor
		err = h.repo.CreateDoctor(userRequest.DoctorDetails)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not created"})
			return
		}
	}

	c.JSON(http.StatusCreated, userRequest)
}

// Login godoc
// @Summary      Authentication
// @Description  Authenticate user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body  models.LoginRequest  true  "Login request object"
// @Success      201  {object}  models.LoginRequest
// @Failure      400  {object}  map[string]string
// @Router       /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var loginRequest models.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by iin
	user, err := h.repo.GetUserByIin(loginRequest.Iin)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Verify password
	if err := auth.VerifyPassword(user.Password, loginRequest.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	var u = uint(user.UserId)
	token, err := auth.GenerateToken(u, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Delete a user by its ID
// @Tags         users
// @Produce      json
// @Param        id  path  int  true  "User ID"
// @Success      204
// @Failure      404  {object}  map[string]string
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, err := h.repo.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	switch user.Role {
	case "patient":
		if err := h.repo.DeletePatient(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error UserHandler - DeleteUser - DeletePatient": err.Error()})
		}
	case "doctor":
		if err := h.repo.DeleteDoctor(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error UserHandler - DeleteUser - DeleteDoctor": err.Error()})
		}
	}

	if err := h.repo.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error UserHandler - DeleteUser - DeleteUser": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

//// GetUserByIIN godoc
//// @Summary      Get a user by IIN
//// @Description  Fetch a user by its IIN
//// @Tags         users
//// @Produce      json
//// @Param        iin  path  string  true  "User IIN"
//// @Success      200  {object}  models.User
//// @Failure      404  {object}  map[string]string
//// @Router       /users/iin/{iin} [get]
//func (h *UserHandler) GetUserByIIN(c *gin.Context) {
//	iin, err := strconv.Atoi(c.Param("iin"))
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IIN"})
//		return
//	}
//
//	user, err := h.repo.GetUserByID(iin)
//	if err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
//		return
//	}
//	c.JSON(http.StatusOK, user)
//}

//// CreateBook godoc
//// @Summary      Create a new book
//// @Description  Create a new book with the provided details
//// @Tags         books
//// @Accept       json
//// @Produce      json
//// @Param        book  body  models.Book  true  "Book object"
//// @Success      201  {object}  models.Book
//// @Failure      400  {object}  map[string]string
//// @Router       /books [post]
//func (h *BookHandler) CreateBook(c *gin.Context) {
//	var book models.Book
//	if err := c.ShouldBindJSON(&book); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	if err := h.repo.CreateBook(&book); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusCreated, book)
//}
//
//// UpdateBook godoc
//// @Summary      Update a book
//// @Description  Update a book by its ID
//// @Tags         books
//// @Accept       json
//// @Produce      json
//// @Param        id    path  int         true  "Book ID"
//// @Param        book  body  models.Book  true  "Updated book object"
//// @Success      200  {object}  models.Book
//// @Failure      400  {object}  map[string]string
//// @Failure      404  {object}  map[string]string
//// @Router       /books/{id} [put]
//func (h *BookHandler) UpdateBook(c *gin.Context) {
//	id, err := strconv.Atoi(c.Param("id"))
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
//		return
//	}
//
//	var book models.Book
//	if err := c.ShouldBindJSON(&book); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	book.ID = id
//
//	if err := h.repo.UpdateBook(&book); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, book)
//}
