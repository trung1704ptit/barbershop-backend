package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/models"
	"gorm.io/gorm"
)

type PointController struct {
	DB *gorm.DB
}

func NewPointController(DB *gorm.DB) PointController {
	return PointController{DB}
}

func (pc *PointController) CreatePoint(ctx *gin.Context) {
	var payload *models.CreatePointRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	now := time.Now()
	newPoint := models.Point{
		User:        payload.User,
		Points:      payload.Points,
		Description: payload.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	result := pc.DB.Create(&newPoint)
	if result.Error != nil {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Point for that user already exists"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newPoint})
}

func (pc *PointController) UpdatePoint(ctx *gin.Context) {
	pointId := ctx.Param("pointId")

	var payload *models.UpdatePointRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var updatedPoint models.Point
	result := pc.DB.First(&updatedPoint, "id = ?", pointId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No point with that title exists"})
		return
	}
	now := time.Now()
	pointToUpdate := models.Point{
		User:        payload.User,
		Points:      payload.Points,
		Description: payload.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	pc.DB.Model(&updatedPoint).Updates(pointToUpdate)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedPoint})
}

func (pc *PointController) FindPointById(ctx *gin.Context) {
	pointId := ctx.Param("pointId")
	var point models.Point
	result := pc.DB.First(&point, "id = ?", pointId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Point not exists"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": point})
}

func (pc *PointController) FindPoints(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)

	offset := (intPage - 1) * intLimit

	var points []models.Point
	results := pc.DB.Limit(intLimit).Offset(offset).Find(&points)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": results.Error})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(points), "data": points})
}

func (pc *PointController) DeletePoint(ctx *gin.Context) {
	pointId := ctx.Param("pointId")
	result := pc.DB.Delete(&models.Point{}, "id = ?", pointId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Point not exists"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
