package controllers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"barbershop-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type ProductController struct {
	DB *gorm.DB
}

func NewProductController(DB *gorm.DB) ProductController {
	return ProductController{DB}
}

func (pc *ProductController) CreateProduct(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(10 << 20) // Max file size 20MB
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Error parsing form data"})
		return
	}

	title := ctx.PostForm("title")
	description := ctx.PostForm("description")
	priceStr := ctx.PostForm("price")
	priceText := ctx.PostForm("price_text")

	price, err := strconv.ParseFloat(priceStr, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid price"})
		return
	}

	previewImg, previewImgHeader, err := ctx.Request.FormFile("preview_image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving preview image"})
		return
	}

	defer previewImg.Close()

	previewImgFilename := uuid.New().String() + filepath.Ext(previewImgHeader.Filename)
	previewImgPath := filepath.Join("uploads", previewImgFilename)
	previewImgFile, err := os.Create(previewImgPath)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Error when save preview image"})
		return
	}

	defer previewImgFile.Close()

	_, err = io.Copy(previewImgFile, previewImg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Error when save preview image"})
		return
	}

	// Get the uploaded images (list of images)
	images := ctx.Request.MultipartForm.File["images"]
	imagePaths := make([]string, len(images))

	for i, img := range images {
		file, err := img.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving image"})
			return
		}
		defer file.Close()

		// Generate a unique filename for each uploaded image
		imageFilename := uuid.New().String() + filepath.Ext(img.Filename)

		// Save the image to a specified location on the server
		imagePath := filepath.Join("uploads", imageFilename)
		imageFile, err := os.Create(imagePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving image"})
			return
		}
		defer imageFile.Close()
		_, err = io.Copy(imageFile, file)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving image"})
			return
		}

		imagePaths[i] = imagePath
	}

	now := time.Now()
	newProduct := models.Product{
		Title:        title,
		Slug:         slug.Make(title),
		Description:  description,
		Price:        price,
		PriceText:    priceText,
		PreviewImage: previewImgPath,
		Images:       imagePaths,
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
		Slug:         slug.Make(payload.Title),
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
