package handlers

import (
	"net/http"
	"strconv"

	//"strconv"
	//
	//"diploma/internal/models"
	"diploma/internal/repositories"
	"github.com/gin-gonic/gin"
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

	if err := h.repo.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
