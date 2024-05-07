package routes

import (
	"barbershop-backend/middleware"
	"barbershop-backend/reminder"

	"github.com/gin-gonic/gin"
)

type RemindRouteController struct {
	remindController reminder.RemindController
}

func NewRouteRemindController(RemindController reminder.RemindController) RemindRouteController {
	return RemindRouteController{RemindController}
}

func (rc *RemindRouteController) RemindRoute(rg *gin.RouterGroup) {

	router := rg.Group("reminds")
	router.Use(middleware.DeserializeUser())
	router.POST("/birthday", rc.remindController.AdminSendBirthdayReminder)
}
