package routes

import (
	"diploma/internal/api/handlers"
	"diploma/internal/auth"
	"diploma/internal/config"
	"diploma/internal/repositories"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to the database
	db, err := repositories.ConnectDB(cfg)
	if err != nil {
		panic(err)
	}

	// Initialize repository and handlers
	userRepo := repositories.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	patientRepo := repositories.NewPatientRepository(db)
	patientHandler := handlers.NewPatientHandler(patientRepo)

	notificactionRepo := repositories.NewNotificationRepository(db)
	notificactionHandler := handlers.NewNotificationHandler(notificactionRepo)

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	v1 := router.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", userHandler.CreateUser)
			authGroup.POST("/login", userHandler.Login)
		}

		usersGroup := v1.Group("/users")
		{
			usersGroup.GET("/", userHandler.GetUsers)
			usersGroup.GET("/:id", userHandler.GetUserByID)
			usersGroup.POST("/", userHandler.CreateUser)
			usersGroup.DELETE("/:id", userHandler.DeleteUser)
			//usersGroup.GET("/iin/:iin", userHandler.GetUserByIIN)
		}

		patientsGroup := v1.Group("/patients")
		{
			patientsGroup.GET("/", patientHandler.GetPatients)
			patientsGroup.GET("/:id", patientHandler.GetPatientByID)
		}

		notificationsGroup := v1.Group("/notifications")
		notificationsGroup.Use(auth.AuthMiddleware())
		{
			notificationsGroup.GET("/:id", notificactionHandler.GetNotificationByUserID)
			notificationsGroup.DELETE("/:id", notificactionHandler.DeleteNotification)
			notificationsGroup.POST("/", notificactionHandler.CreateNotification)
		}

	}

	return router
}
