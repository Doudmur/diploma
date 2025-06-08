package handlers

import (
	"bytes"
	"diploma/internal/auth"
	"diploma/internal/models"
	"diploma/internal/repositories"
	"diploma/internal/scripts"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
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

	// Check if user exists
	if user, _ := h.repo.GetUserByIin(userRequest.Iin); user != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this IIN already exists"})
		return
	}

	// Check if user exists
	if user, _ := h.repo.GetUserByEmail(userRequest.Email); user != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(otp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	userRequest.Password = hashedPassword

	// Create User and related records in a transaction
	err = h.repo.CreateUser(&userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully. Please upload your photo for face verification.",
		"user_id": userRequest.UserId,
	})
}

// UploadPhoto godoc
// @Summary      Upload user photo for face verification
// @Description  Upload and verify user's photo
// @Tags         auth
// @Accept       multipart/form-data
// @Produce      json
// @Param        user_id formData int true "User ID"
// @Param        photo formData file true "User Photo"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /auth/upload-photo [post]
func (h *UserHandler) UploadPhoto(c *gin.Context) {
	userID, err := strconv.Atoi(c.PostForm("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get photo file
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Photo is required"})
		return
	}

	// Open the file
	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open photo"})
		return
	}
	defer openedFile.Close()

	// Read file content
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, openedFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read photo"})
		return
	}
	photoBytes := buf.Bytes()

	// Send to Python face detection
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("file", file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare photo for face detection"})
		return
	}
	fw.Write(photoBytes)
	w.Close()

	req, err := http.NewRequest("POST", "http://134.122.84.85:8000/detect_face/", &b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create face detection request"})
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate face"})
		return
	}
	defer resp.Body.Close()

	var detectResp models.DetectResponse
	if err := json.NewDecoder(resp.Body).Decode(&detectResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode face detection response"})
		return
	}

	if !detectResp.HasFace || detectResp.FaceCount != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image must contain exactly one face"})
		return
	}

	// Update user's photo in database
	err = h.repo.UpdateUserPhoto(userID, photoBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save photo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo uploaded and verified successfully"})
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

	user, err := h.repo.GetUserByIin(loginRequest.Iin)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !auth.CheckPasswordHash(loginRequest.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if password change is required
	if !user.PasswordChanged {
		c.JSON(http.StatusOK, gin.H{
			"message":                 "Password change required",
			"require_password_change": true,
			"token":                   nil,
		})
		return
	}

	token, err := auth.GenerateToken(uint(user.UserId), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// ChangePassword godoc
// @Summary      Change password
// @Description  Change user's password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body  models.ChangePasswordRequest  true  "Change Password Request"
// @Param        Authorization header string true "Bearer"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /users/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var request models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repo.GetUserByIin(request.Iin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !auth.CheckPasswordHash(request.OldPassword, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid old password"})
		return
	}

	if request.NewPassword != request.VerifyPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New passwords do not match"})
		return
	}

	hashedPassword, err := auth.HashPassword(request.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if err := h.repo.UpdatePassword(request.Iin, hashedPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	if err := h.repo.SetPasswordChanged(request.Iin, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password_changed flag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
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

// ForgotPassword godoc
// @Summary      Request password reset
// @Description  Send OTP to user's email for password reset
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body  models.ForgotPasswordRequest  true  "Forgot Password Request"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /users/forgot-password [post]
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var request models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.repo.GetUserByIin(request.Iin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	otp := generateOTP()
	expiresAt := time.Now().Add(5 * time.Minute)

	otpVerification := models.OTPVerification{
		Iin:       request.Iin,
		OTP:       otp,
		ExpiresAt: expiresAt,
	}

	if err := h.repo.CreateOTPVerification(&otpVerification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OTP verification"})
		return
	}

	// TODO: Send OTP via email
	// For now, just return the OTP in the response
	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent to email",
		"otp":     otp, // Remove this in production
	})
}

// VerifyOTP godoc
// @Summary      Verify OTP
// @Description  Verify OTP and set password_changed flag to false
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body  models.VerifyOTPRequest  true  "Verify OTP Request"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /users/verify-otp [post]
func (h *UserHandler) VerifyOTP(c *gin.Context) {
	var request models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	otpVerification, err := h.repo.GetOTPVerification(request.Iin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "OTP verification not found"})
		return
	}

	if otpVerification.OTP != request.OTP {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	if time.Now().After(otpVerification.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP has expired"})
		return
	}

	if err := h.repo.SetPasswordChanged(request.Iin, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password_changed flag"})
		return
	}

	if err := h.repo.DeleteOTPVerification(request.Iin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete OTP verification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}

func generateOTP() string {
	// Generate a 6-digit OTP
	otp := rand.Intn(900000) + 100000
	return fmt.Sprintf("%06d", otp)
}

// GetUserInfoByIIN godoc
// @Summary      Get detailed user information by IIN
// @Description  Get user information including role-specific details (doctor or patient)
// @Tags         users
// @Produce      json
// @Param        iin  path  string  true  "User IIN"
// @Success      200  {object}  models.UserInfoResponse
// @Failure      404  {object}  map[string]string
// @Router       /users/info/{iin} [get]
func (h *UserHandler) GetUserInfoByIIN(c *gin.Context) {
	iin := c.Param("iin")
	if iin == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IIN is required"})
		return
	}

	userInfo, err := h.repo.GetUserInfoByIIN(iin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "Detailed": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userInfo)
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
//
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
