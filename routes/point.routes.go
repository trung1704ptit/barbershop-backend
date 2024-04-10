package routes

import (
	"barbershop-backend/controllers"
	"barbershop-backend/middleware"

	"github.com/gin-gonic/gin"
)

type PointRouteController struct {
	pointController controllers.PointController
}

func NewRoutePointController(pointController controllers.PointController) PointRouteController {
	return PointRouteController{pointController}
}

func (pc *PointRouteController) PointRoute(rg *gin.RouterGroup) {
	router := rg.Group("points")
	router.POST("/", middleware.DeserializeUser(), pc.pointController.CreatePoint)
	router.GET("/", pc.pointController.FindPoints)
	router.GET("/:pointId", pc.pointController.FindPointById)
	router.PUT("/:pointId", middleware.DeserializeUser(), pc.pointController.UpdatePoint)
	router.DELETE("/:pointId", middleware.DeserializeUser(), pc.pointController.DeletePoint)
}
