package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/controllers"
	"github.com/wpcodevo/golang-gorm-postgres/middleware"
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
