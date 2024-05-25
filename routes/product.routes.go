package routes

import (
	"barbershop-backend/controllers"
	"barbershop-backend/middleware"

	"github.com/gin-gonic/gin"
)

type ProductRouteController struct {
	productController controllers.ProductController
}

func NewRouteProductController(productController controllers.ProductController) ProductRouteController {
	return ProductRouteController{productController}
}

func (pc *ProductRouteController) PostRoute(rg *gin.RouterGroup) {

	router := rg.Group("products")
	router.POST("", middleware.DeserializeUser(), pc.productController.CreateProduct)
	router.GET("", pc.productController.FindProducts)
	router.PUT(":productId", middleware.DeserializeUser(), pc.productController.UpdateProduct)
	router.GET(":productId", pc.productController.FindProductById)
	router.DELETE(":productId", middleware.DeserializeUser(), pc.productController.DeleteProduct)
}
