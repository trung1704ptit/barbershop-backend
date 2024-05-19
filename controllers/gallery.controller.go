package controllers

import (
	"barbershop-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GalleryController struct {
	DB *gorm.DB
}

func NewGalleryController(DB *gorm.DB) GalleryController {
	return GalleryController{DB}
}

// CreateGallery creates a new gallery item
func (gc *GalleryController) CreateGallery(ctx *gin.Context) {
	var gallery models.Gallery
	if err := ctx.ShouldBindJSON(&gallery); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	gallery.ID = uuid.New()
	if err := gc.DB.Create(&gallery).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gallery})
}

// GetGalleries retrieves all gallery items
func (gc *GalleryController) GetGalleries(ctx *gin.Context) {
	var galleries []models.Gallery
	if err := gc.DB.Find(&galleries).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": galleries})
}

// UpdateGallery updates a gallery item by ID
func (gc *GalleryController) UpdateGallery(ctx *gin.Context) {
	id := ctx.Param("id")
	var gallery models.Gallery
	if err := gc.DB.First(&gallery, "id = ?", id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Gallery not found"})
		return
	}
	if err := ctx.ShouldBindJSON(&gallery); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	if err := gc.DB.Save(&gallery).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gallery})
}
