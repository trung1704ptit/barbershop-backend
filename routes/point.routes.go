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
	router.POST("", middleware.DeserializeUser(), pc.pointController.CreatePoint)
	router.GET(":phone", pc.pointController.FindPointsByPhone)
	router.GET("", pc.pointController.FindPoints)
	router.PUT(":pointId", middleware.DeserializeUser(), pc.pointController.UpdatePoint)
	router.DELETE(":pointId", middleware.DeserializeUser(), pc.pointController.DeleteOnePoint)
	router.DELETE("/history/:userId", middleware.DeserializeUser(), pc.pointController.DeleteAllPointsByUserId)
}
