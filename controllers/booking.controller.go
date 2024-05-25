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
func (bc *BookingController) UpdateBooking(ctx *gin.Context) {
	var payload models.UpdateBookingRequest

	// Bind the JSON payload to the struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Check if the booking exists
	var booking models.Booking
	if err := bc.DB.First(&booking, "id = ?", payload.ID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Booking not found"})
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

	// Start a transaction
	tx := bc.DB.Begin()

	// Update the booking details
	booking.BarberID = payload.BarberID
	booking.GuestID = payload.GuestID
	booking.Status = payload.Status
	booking.BookingTime = payload.BookingTime
	booking.UpdatedAt = time.Now()

	// Save the booking
	if err := tx.Save(&booking).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Update the services associated with the booking
	var services []models.Service
	if err := tx.Where("id IN ?", payload.ServiceIDs).Find(&services).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Some services not found"})
		return
	}

	// Clear existing services and associate new ones
	if err := tx.Model(&booking).Association("Services").Clear(); err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	if err := tx.Model(&booking).Association("Services").Append(services); err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Commit the transaction
	tx.Commit()

	// Preload the associated Barber and Services
	if err := bc.DB.Preload("Barber").Preload("Guest").Preload("Services").First(&booking, "id = ?", booking.ID).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": booking})
}

func (pc *BookingController) FindBookings(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "1000")
	var barberId = ctx.DefaultQuery("barber_id", "*")
	var status = ctx.DefaultQuery("status", "*")
	var month = ctx.DefaultQuery("month", "*")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	intMonth, _ := strconv.Atoi(month)
	offset := (intPage - 1) * intLimit

	var bookings []models.Booking
	query := pc.DB.Preload("Barber").Preload("Guest").Preload("Services").Limit(intLimit).Offset(offset)

	if barberId != "*" {
		query = query.Where("barber_id = ?", barberId)
	}

	if status != "*" {
		query = query.Where("status = ?", status)
	}

	if month != "*" {
		startDate := time.Date(time.Now().Year(), time.Month(intMonth), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0)
		query = query.Where("booking_time >= ? AND booking_time < ?", startDate, endDate)
	}

	results := query.Find(&bookings)
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

func (bc *BookingController) DeleteBooking(ctx *gin.Context) {
	bookingID := ctx.Param("bookingId")

	var booking models.Booking
	if err := bc.DB.Preload("Services").First(&booking, "id = ?", bookingID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Booking not found"})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		}
		return
	}

	// Start a transaction
	tx := bc.DB.Begin()

	// Remove the association with services
	if err := tx.Model(&booking).Association("Services").Clear(); err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not remove service associations"})
		return
	}

	// Delete the booking
	if err := tx.Delete(&booking).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not delete booking"})
		return
	}

	// Commit the transaction
	tx.Commit()

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Booking deleted successfully"})
}
func (bc *BookingController) DeleteBookingsByGuestID(ctx *gin.Context) {
	guestID := ctx.Param("guestId")

	// Start a transaction
	tx := bc.DB.Begin()

	// Find all bookings by guestID
	var bookings []models.Booking
	if err := tx.Preload("Services").Where("guest_id = ?", guestID).Find(&bookings).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not find bookings by guest ID"})
		return
	}

	// Remove associations with services for each booking
	for _, booking := range bookings {
		if err := tx.Model(&booking).Association("Services").Clear(); err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not remove service associations"})
			return
		}
	}

	// Delete all bookings by guestID
	if err := tx.Where("guest_id = ?", guestID).Delete(&models.Booking{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not delete bookings"})
		return
	}

	// Commit the transaction
	tx.Commit()

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Bookings deleted successfully"})
}
