package routes

import (
	"barbershop-backend/controllers"
	"barbershop-backend/middleware"

	"github.com/gin-gonic/gin"
)

type FileRouteController struct {
	fileController controllers.FileController
}

func NewRouteFileController(fileController controllers.FileController) FileRouteController {
	return FileRouteController{fileController}
}

func (fc *FileRouteController) FileRoute(rg *gin.RouterGroup) {
	router := rg.Group("files")
	router.POST("/single", middleware.DeserializeUser(), fc.fileController.Upload)
	router.POST("/multiple", middleware.DeserializeUser(), fc.fileController.MultiUpload)

}
