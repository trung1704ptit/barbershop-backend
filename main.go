package main

import (
	"log"
	"net/http"
	"time"

	"barbershop-backend/controllers"
	"barbershop-backend/initializers"
	"barbershop-backend/reminder"
	"barbershop-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	server = gin.Default()
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

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	PostRouteController.PostRoute(router)
	RemindRouteController.RemindRoute(router)
	ServiceRouteController.ServiceRoute(router)
	PointRouteController.PointRoute(router)
	log.Fatal(server.Run(":" + config.ServerPort))

	go sleepUntilNext10AM()

	// Other code can continue executing here
	// For example:
	log.Println("This code will execute while waiting for 10 AM...")

	// Block main goroutine to keep the program running
	select {}
}

func sleepUntilNext10AM() {
	startOfDay := time.Now().Truncate(24 * time.Hour)
	desiredTime := startOfDay.Add(10 * time.Hour)

	// Calculate the duration until the next 10 AM
	durationUntilNext10AM := desiredTime.Sub(time.Now())
	if durationUntilNext10AM < 0 {
		durationUntilNext10AM += 24 * time.Hour
	}

	// Sleep until the next 10 AM
	time.Sleep(durationUntilNext10AM)

	// Start the loop to check and send birthday reminds every 24 hours at 10 AM
	for {
		RemindController.CheckAndSendBirthdayReminders()
		time.Sleep(24 * time.Hour) // Check every 24 hours
	}
}
