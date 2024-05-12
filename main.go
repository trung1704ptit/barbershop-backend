package main

import (
	"log"
	"net/http"

	"barbershop-backend/controllers"
	"barbershop-backend/initializers"
	"barbershop-backend/reminder"
	"barbershop-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	PostController      controllers.PostController
	PostRouteController routes.PostRouteController

	ServiceController      controllers.ServiceController
	ServiceRouteController routes.ServiceRouteController

	PointController      controllers.PointController
	PointRouteController routes.PointRouteController

	RemindController      reminder.RemindController
	RemindRouteController routes.RemindRouteController

	FileController      controllers.FileController
	FileRouteController routes.FileRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(initializers.DB)
	UserRouteController = routes.NewRouteUserController(UserController)

	PostController = controllers.NewPostController(initializers.DB)
	PostRouteController = routes.NewRoutePostController(PostController)

	ServiceController = controllers.NewServiceController(initializers.DB)
	ServiceRouteController = routes.NewRouteServiceController(ServiceController)

	PointController = controllers.NewPointController(initializers.DB)
	PointRouteController = routes.NewRoutePointController(PointController)

	RemindController = reminder.NewRemindController(initializers.DB, &UserController)
	RemindRouteController = routes.NewRouteRemindController(RemindController)

	FileController = controllers.NewFileController(initializers.DB)
	FileRouteController = routes.NewRouteFileController(FileController)

	server = gin.Default()

	server.MaxMultipartMemory = 100 << 20 // 50MB
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", "http://localhost:3000", config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Golang with Gorm and Postgres"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	router.Static("/uploads", "./uploads")

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	PostRouteController.PostRoute(router)
	RemindRouteController.RemindRoute(router)
	ServiceRouteController.ServiceRoute(router)
	PointRouteController.PointRoute(router)
	FileRouteController.FileRoute(router)

	// Create a new cron job scheduler
	c := cron.New()

	// Schedule the email to be sent every day at 9:00 AM
	_, cronError := c.AddFunc("0 8 * * *", sendEmail)
	if cronError != nil {
		log.Fatal("Error adding cron job:", cronError)
	}

	// Start the cron scheduler
	c.Start()

	log.Fatal(server.Run(":" + config.ServerPort))

	// Keep the program running to allow cron jobs to execute
	select {}
}

func sendEmail() {
	log.Println("Sending email...")

	// Send the email here
	RemindController.CheckAndSendBirthdayReminders()
	log.Println("Email sent!")
}
