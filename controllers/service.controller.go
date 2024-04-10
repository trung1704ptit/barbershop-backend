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

type ServiceController struct {
	DB *gorm.DB
}

func NewServiceController(DB *gorm.DB) ServiceController {
	return ServiceController{DB}
}

func (sc *ServiceController) CreateService(ctx *gin.Context) {
	var payload *models.CreateServiceRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	newService := models.Service{
		Name:        payload.Name,
		Image:       payload.Image,
		Price:       payload.Price,
		PriceText:   payload.PriceText,
		Description: payload.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	result := sc.DB.Create(&newService)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "service with that name already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": result.Error.Error()})
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newService})
}

func (sc *ServiceController) UpdateService(ctx *gin.Context) {
	serviceId := ctx.Param("serviceId")

	var payload *models.UpdateService
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var updatedService models.Service
	result := sc.DB.First(&updatedService, "id = ?", serviceId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No service with that name exists"})
		return
	}
	now := time.Now()
	serviceToUpdate := models.Service{
		Name:        payload.Name,
		Image:       payload.Image,
		Price:       payload.Price,
		PriceText:   payload.PriceText,
		Description: payload.Description,
		UpdatedAt:   now,
		CreatedAt:   payload.CreatedAt,
	}
	sc.DB.Model(&updatedService).Updates(serviceToUpdate)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedService})
}

func (sc *ServiceController) FindServiceById(ctx *gin.Context) {
	serviceId := ctx.Param("serviceId")
	var service models.Service
	result := sc.DB.First(&service, "id = ?", serviceId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No Service with that name exists"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": service})
}

func (sc *ServiceController) FindServices(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var services []models.Service

	results := sc.DB.Limit(intLimit).Offset(offset).Find(&services)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": results.Error})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(services), "data": services})
}

func (sc *ServiceController) DeleteService(ctx *gin.Context) {
	serviceId := ctx.Param("serviceId")
	result := sc.DB.Delete(&models.Service{}, "id = ?", serviceId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No service with that name exists"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
