package router

import (
	controllers "SanjaySinghRajpoot/ticketBooking/controller"
	"SanjaySinghRajpoot/ticketBooking/middleware"

	"github.com/gin-gonic/gin"
)

func GetRoute(r *gin.Engine) {

	// admin routes
	// add trains
	r.POST("/api/train/add", controllers.SaveTrain)

	// User routes
	r.POST("/api/signup", controllers.Signup)
	r.POST("/api/login", controllers.Login)

	r.Use(middleware.RequireAuth)
	r.POST("/api/logout", controllers.Logout)
	r.GET("/api/trains", controllers.CheckAvailableTrains)
	r.GET("/api/trains/seats", controllers.CheckAvailableSeats)

	r.POST("/api/booking", controllers.Booking)
	r.GET("/api/booking", controllers.GetBookingData)
}
