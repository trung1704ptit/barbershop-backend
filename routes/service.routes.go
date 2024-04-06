package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/controllers"
	"github.com/wpcodevo/golang-gorm-postgres/middleware"
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
	router.GET("/", sc.serviceController.FindServices)
	router.GET("/:serviceId", sc.serviceController.FindServiceById)
	router.PUT("/:serviceId", middleware.DeserializeUser(), sc.serviceController.UpdateService)
	router.DELETE("/:serviceId", middleware.DeserializeUser(), sc.serviceController.DeleteService)
}
