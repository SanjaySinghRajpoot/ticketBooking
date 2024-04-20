package router

import (
	controllers "SanjaySinghRajpoot/ticketBooking/controller"
	"SanjaySinghRajpoot/ticketBooking/middleware"

	"github.com/gin-gonic/gin"
)

func GetRoute(r *gin.Engine) {

	// admin routes
	r.POST("/api/train", controllers.SaveTrain)
	r.DELETE("/api/train", controllers.DeleteTrain)
	r.PATCH("/api/train", controllers.UpdateTrain)

	// User routes
	r.POST("/api/signup", controllers.Signup)
	r.POST("/api/login", controllers.Login)
	r.POST("/api/logout", controllers.Logout)

	// Trains Routes
	r.Use(middleware.RequireAuth)
	r.GET("/api/trains", controllers.CheckAvailableTrains)
	r.GET("/api/trains/seats", controllers.CheckAvailableSeats)

	// Booking Routes
	r.POST("/api/booking", controllers.Booking)
	r.GET("/api/booking", controllers.GetBookingData)
	r.POST("/api/booking/cancel", controllers.CancelBooking)

}
