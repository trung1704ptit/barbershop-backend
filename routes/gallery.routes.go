package routes

import (
	"barbershop-backend/controllers"

	"github.com/gin-gonic/gin"
)

type GalleryRouteController struct {
	galleryController controllers.GalleryController
}

func NewRouteGalleryController(galleryController controllers.GalleryController) GalleryRouteController {
	return GalleryRouteController{galleryController}
}

func (fc *GalleryRouteController) GalleryRoute(rg *gin.RouterGroup) {
	router := rg.Group("galleries")
	router.POST("", fc.galleryController.CreateGallery)
	router.GET("", fc.galleryController.GetGalleries)
	router.PUT("/:id", fc.galleryController.UpdateGallery)
}
