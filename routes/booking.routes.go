package routes

import (
	"barbershop-backend/controllers"
	"barbershop-backend/middleware"

	"github.com/gin-gonic/gin"
)

type BookingRouteController struct {
	bookingController controllers.BookingController
}

func NewRouteBookingController(bookingController controllers.BookingController) BookingRouteController {
	return BookingRouteController{bookingController}
}

func (bc *BookingRouteController) BookingRoute(rg *gin.RouterGroup) {
	router := rg.Group("bookings")
	router.POST("", bc.bookingController.CreateBooking)
	router.GET("", middleware.DeserializeUser(), bc.bookingController.FindBookings)
	router.PUT(":bookingId", middleware.DeserializeUser(), bc.bookingController.UpdateBooking)
	router.GET(":bookingId", bc.bookingController.FindBookingById)
	router.DELETE(":bookingId", middleware.DeserializeUser(), bc.bookingController.DeleteBooking)
}
