package routes

import (
	"barbershop-backend/controllers"
	"barbershop-backend/middleware"

	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {

	router := rg.Group("users")
	router.GET("", uc.userController.FindUsers)
	router.GET("/me", middleware.DeserializeUser(), uc.userController.GetMe)
	router.GET("/:phone", uc.userController.GetUserByPhone)
	router.PUT("/:userId", middleware.DeserializeUser(), uc.userController.UpdateUser)
	router.DELETE("/:userId", middleware.DeserializeUser(), uc.userController.DeleteUser)
}
