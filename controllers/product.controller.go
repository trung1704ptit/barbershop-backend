package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"barbershop-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController struct {
	DB *gorm.DB
}

func NewProductController(DB *gorm.DB) ProductController {
	return ProductController{DB}
}

func (pc *ProductController) CreateProduct(ctx *gin.Context) {
	var payload *models.CreateProductRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	newProduct := models.Product{
		Title:        payload.Title,
		Slug:         payload.Slug,
		Description:  payload.Description,
		Price:        payload.Price,
		PriceText:    payload.PriceText,
		PreviewImage: payload.PreviewImage,
		Images:       payload.Images,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result := pc.DB.Create(&newProduct)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Product with that title already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newProduct})
}

func (pc *ProductController) UpdateProduct(ctx *gin.Context) {
	productId := ctx.Param("productId")

	var payload *models.UpdateProduct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var updatedProduct models.Product
	result := pc.DB.First(&updatedProduct, "id = ?", productId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No product with that title exists"})
		return
	}
	now := time.Now()
	productToUpdate := models.Product{
		Title:        payload.Title,
		Slug:         payload.Slug,
		Description:  payload.Description,
		Price:        payload.Price,
		PriceText:    payload.PriceText,
		PreviewImage: payload.PreviewImage,
		Images:       payload.Images,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	pc.DB.Model(&updatedProduct).Updates(productToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedProduct})
}

func (pc *ProductController) FindProductById(ctx *gin.Context) {
	productId := ctx.Param("productId")

	var product models.Product
	result := pc.DB.First(&product, "id = ?", productId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No product with that title exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": product})
}

func (pc *ProductController) FindProducts(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var products []models.Product
	results := pc.DB.Limit(intLimit).Offset(offset).Find(&products)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(products), "data": products})
}

func (pc *ProductController) DeleteProduct(ctx *gin.Context) {
	productId := ctx.Param("productId")

	result := pc.DB.Delete(&models.Product{}, "id = ?", productId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No product with that title exists"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
