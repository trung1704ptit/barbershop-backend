package controllers

import (
	"net/http"
	"strconv"
	"time"

	"barbershop-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingController struct {
	DB *gorm.DB
}

func NewBookingController(DB *gorm.DB) BookingController {
	return BookingController{DB}
}

func (pc *BookingController) CreateBooking(ctx *gin.Context) {
	var payload *models.CreateBookingRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()
	newBooking := models.Booking{
		BarberID:  payload.BarberID,
		GuestID:   payload.GuestID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := pc.DB.Create(&newBooking)
	if result.Error != nil {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	pc.DB.Preload("Barber").Preload("Guest").First(&newBooking, "id = ?", newBooking.ID)

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newBooking})
}

func (pc *BookingController) FindBookings(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "1000")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var bookings []models.Booking
	results := pc.DB.Preload("Barber").Preload("Guest").Limit(intLimit).Offset(offset).Find(&bookings)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(bookings), "data": bookings})
}

func (pc *BookingController) FindBookingById(ctx *gin.Context) {
	bookingId := ctx.Param("bookingId")

	var booking models.Booking
	result := pc.DB.First(&booking, "id = ?", bookingId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No booking with that id exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": booking})
}

func (pc *BookingController) DeleteBooking(ctx *gin.Context) {
	bookingId := ctx.Param("bookingId")

	result := pc.DB.Delete(&models.Booking{}, "id = ?", bookingId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No booking with that id exists"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
