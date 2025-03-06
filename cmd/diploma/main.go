package main

import (
	_ "diploma/docs" // Import the generated docs
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"strconv"
)

// Product represents a product entity
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var products = []Product{
	{ID: 1, Name: "Laptop", Price: 999.99},
	{ID: 2, Name: "Smartphone", Price: 499.99},
}

// @title My API
// @version 1.0
// @description This is a sample API with Swagger documentation.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	r := gin.Default()

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/hello", helloHandler)
		v1.GET("/products", getProducts)
		v1.GET("/products/:id", getProduct)
		v1.POST("/products", createProduct)
		v1.PUT("/products/:id", updateProduct)
		v1.DELETE("/products/:id", deleteProduct)
	}

	r.Run(":8080")
}

// @Summary Say hello
// @Description Get a friendly greeting
// @Produce json
// @Success 200 {string} string "Hello, World!"
// @Router /hello [get]
func helloHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}

// @Summary Get all products
// @Description Get a list of all products
// @Produce json
// @Success 200 {array} Product
// @Router /products [get]
func getProducts(c *gin.Context) {
	c.JSON(http.StatusOK, products)
}

// @Summary Get a product by ID
// @Description Get a product by its ID
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} Product
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func getProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	for _, product := range products {
		if product.ID == id {
			c.JSON(http.StatusOK, product)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}

// @Summary Create a new product
// @Description Create a new product with the input payload
// @Accept json
// @Produce json
// @Param product body Product true "Product object"
// @Success 201 {object} Product
// @Failure 400 {object} map[string]string
// @Router /products [post]
func createProduct(c *gin.Context) {
	var newProduct Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Assign a new ID
	newProduct.ID = len(products) + 1
	products = append(products, newProduct)

	c.JSON(http.StatusCreated, newProduct)
}

// @Summary Update a product
// @Description Update a product by its ID
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body Product true "Product object"
// @Success 200 {object} Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [put]
func updateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var updatedProduct Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	for i, product := range products {
		if product.ID == id {
			updatedProduct.ID = id
			products[i] = updatedProduct
			c.JSON(http.StatusOK, updatedProduct)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}

// @Summary Delete a product
// @Description Delete a product by its ID
// @Produce json
// @Param id path int true "Product ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /products/{id} [delete]
func deleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	for i, product := range products {
		if product.ID == id {
			products = append(products[:i], products[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}
