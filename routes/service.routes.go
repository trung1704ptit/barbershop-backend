package routes

import (
	"barbershop-backend/controllers"
	"barbershop-backend/middleware"

	"github.com/gin-gonic/gin"
)

type ServiceRouteController struct {
	serviceController controllers.ServiceController
}

func NewRouteServiceController(serviceController controllers.ServiceController) ServiceRouteController {
	return ServiceRouteController{serviceController}
}

func (sc *ServiceRouteController) ServiceRoute(rg *gin.RouterGroup) {
	router := rg.Group("services")
	router.POST("/", middleware.DeserializeUser(), sc.serviceController.CreateService)
	router.POST("/user-service", middleware.DeserializeUser(), sc.serviceController.RegisterUserWithServices)
	router.GET("/user-service", middleware.DeserializeUser(), sc.serviceController.FindUserServices)
	router.DELETE("/user-service/:userAndServiceId", middleware.DeserializeUser(), sc.serviceController.DeleteUserWithServices)
	router.POST("/history", middleware.DeserializeUser(), sc.serviceController.AddServiceHistory)
	router.GET("/history", middleware.DeserializeUser(), sc.serviceController.GetServiceHistories)
	router.GET("/", sc.serviceController.FindServices)
	router.GET("/:serviceId", sc.serviceController.FindServiceById)
	router.PUT("/:serviceId", middleware.DeserializeUser(), sc.serviceController.UpdateService)
	router.DELETE("/:serviceId", middleware.DeserializeUser(), sc.serviceController.DeleteService)
}
