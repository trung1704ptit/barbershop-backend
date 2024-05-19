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
		Todos:       payload.Todos,
		Category:    payload.Category,
		ServiceType: payload.ServiceType,
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
		Todos:       payload.Todos,
		Category:    payload.Category,
		ServiceType: payload.ServiceType,
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
	var serviceType = ctx.DefaultQuery("service_type", "*")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var services []models.Service

	query := sc.DB.Limit(intLimit).Offset(offset)

	if serviceType != "*" {
		query = query.Where("service_type = ?", serviceType)
	}

	results := query.Find(&services)
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

func (sc *ServiceController) RegisterUserWithServices(ctx *gin.Context) {
	var payload *models.UserServiceRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Parsing body failed"})
		return
	}

	now := time.Now()
	newUserService := models.UserService{
		UserID:          payload.UserID,
		ServiceID:       payload.ServiceID,
		CreatedAt:       now,
		UpdatedAt:       now,
		UserIDServiceID: payload.UserID.String() + "_" + payload.ServiceID.String(),
	}

	result := sc.DB.Create(&newUserService)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newUserService})
}

func (sc *ServiceController) AddServiceHistory(ctx *gin.Context) {
	var payload *models.UserServiceRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Parsing body failed"})
		return
	}
	var lastRecord *models.ServiceHistory
	var lastCount int

	lastRecordResult := sc.DB.Order("created_at desc").Where("user_id = ? and service_id = ?", payload.UserID, payload.ServiceID).Last(&lastRecord)

	if lastRecordResult.Error != nil {
		lastCount = 0
	} else {
		lastCount = lastRecord.Count
	}

	now := time.Now()
	newServiceHistory := models.ServiceHistory{
		UserID:    payload.UserID,
		ServiceID: payload.ServiceID,
		Count:     lastCount + 1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := sc.DB.Create(&newServiceHistory)
	if result.Error != nil {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newServiceHistory})
}

func (sc *ServiceController) GetServiceHistories(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "100")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var servicesHistory []models.ServiceHistory

	results := sc.DB.Limit(intLimit).Offset(offset).Find(&servicesHistory)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": results.Error})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(servicesHistory), "data": servicesHistory})
}

func (sc *ServiceController) DeleteUserWithServices(ctx *gin.Context) {
	userAndServiceId := ctx.Param("userAndServiceId")

	result := sc.DB.Delete(&models.UserService{}, "user_id_service_id = ?", userAndServiceId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No service with that name exists"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (sc *ServiceController) DeleteServicesHistory(ctx *gin.Context) {
	userId := ctx.Param("userId")

	result := sc.DB.Delete(&models.ServiceHistory{}, "user_id = ?", userId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No service with that name exists"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (sc *ServiceController) FindUserServices(ctx *gin.Context) {
	var userServices []models.UserService

	results := sc.DB.Limit(1000).Offset(0).Find(&userServices)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": results.Error})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(userServices), "data": userServices})
}
