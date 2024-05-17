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

func (bc *BookingController) CreateBooking(ctx *gin.Context) {
	var payload models.CreateBookingRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Check if the barber and guest exist
	var barber models.User
	var guest models.User
	if err := bc.DB.First(&barber, "id = ?", payload.BarberID).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Barber not found"})
		return
	}
	if err := bc.DB.First(&guest, "id = ?", payload.GuestID).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Guest not found"})
		return
	}

	now := time.Now()
	newBooking := models.Booking{
		BarberID:    payload.BarberID,
		GuestID:     payload.GuestID,
		Status:      "open",
		BookingTime: payload.BookingTime,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Start a transaction
	tx := bc.DB.Begin()

	// Create the booking record
	if err := tx.Create(&newBooking).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Find the services
	var services []models.Service
	if err := tx.Where("id IN ?", payload.ServiceIDs).Find(&services).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Some services not found"})
		return
	}

	// Associate the services with the booking
	if err := tx.Model(&newBooking).Association("Services").Append(services); err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Commit the transaction
	tx.Commit()

	// Preload the associated Barber and Services
	if err := bc.DB.Preload("Barber").Preload("Guest").Preload("Services").First(&newBooking, "id = ?", newBooking.ID).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newBooking})
}

func (pc *BookingController) FindBookings(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "1000")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var bookings []models.Booking
	results := pc.DB.Preload("Barber").Preload("Guest").Preload("Services").Limit(intLimit).Offset(offset).Find(&bookings)
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
