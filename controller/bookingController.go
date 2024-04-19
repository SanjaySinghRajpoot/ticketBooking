package controllers

import (
	"SanjaySinghRajpoot/ticketBooking/config"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	mutex sync.Mutex
)

func Booking(c *gin.Context) {
	var booking struct {
		UserID  string `json:"user_id" binding:"required"`
		TrainID string `json:"train_id" binding:"required"`
		Seats   int    `json:"seats" binding:"required"`
	}

	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Lock the mutex to ensure exclusive access to shared resources
	mutex.Lock()
	defer mutex.Unlock()

	// Start a transaction
	tx, err := config.DB.Begin()
	if err != nil {
		fmt.Println("Error starting transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to start transaction",
		})
		return
	}
	defer tx.Rollback() // Rollback transaction on error

	// Lock the selected row in the trains table
	var totalSeats int
	row := tx.QueryRow("SELECT total_seats FROM trains WHERE id = $1 FOR UPDATE", booking.TrainID)
	if err := row.Scan(&totalSeats); err != nil {
		fmt.Println("Error scanning row:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get total seats from the train",
		})
		return
	}

	// Check if the requested seats are available
	if totalSeats < booking.Seats {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": "Requested seats not available",
		})
		return
	}

	// Subtract the booked seats from the total seats
	remainingSeats := totalSeats - booking.Seats

	// Update the total seats in the train
	_, err = tx.Exec("UPDATE trains SET total_seats = $1 WHERE id = $2", remainingSeats, booking.TrainID)
	if err != nil {
		fmt.Println("Error updating total seats:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update total seats in the train",
		})
		return
	}

	// Create a new entry in the bookings table
	_, err = tx.Exec("INSERT INTO bookings (user_id, train_id, seats, status, created_at) VALUES ($1, $2, $3, $4)", booking.UserID, booking.TrainID, booking.Seats, "booked", time.Now())
	if err != nil {
		fmt.Println("Error creating booking:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create booking",
		})
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		fmt.Println("Error committing transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit transaction",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("%d seats booked successfully for train %s", booking.Seats, booking.TrainID),
	})
}

// create endpoint to get the booking details using a given booking id

func GetBookingData(c *gin.Context) {
	bookingID := c.Query("id")
	if bookingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Booking ID is required",
		})
		return
	}

	var (
		id            int
		seats         int
		status        *string
		trainName     string
		departureTime time.Time
		arrivalTime   time.Time
		fromStation   string
		toStation     string
	)

	row := config.DB.QueryRow(`
		SELECT bookings.id, bookings.seats, bookings.status, trains.name, trains.departure_time, trains.arrival_time, trains."from", trains."to"
		FROM bookings
		LEFT JOIN trains ON bookings.train_id = trains.id
		WHERE bookings.id = $1
	`, bookingID)

	err := row.Scan(&id, &seats, &status, &trainName, &departureTime, &arrivalTime, &fromStation, &toStation)
	if err != nil {
		fmt.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch booking data",
		})
		return
	}

	bookingData := gin.H{
		"id":             id,
		"seats":          seats,
		"status":         status,
		"train_name":     trainName,
		"departure_time": departureTime,
		"arrival_time":   arrivalTime,
		"from":           fromStation,
		"to":             toStation,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"bookingData": bookingData,
	})
}

func CancelBooking(c *gin.Context) {

	var booking struct {
		UserID  string `json:"user_id" binding:"required"`
		TrainID string `json:"train_id" binding:"required"`
		Seats   int    `json:"seats" binding:"required"`
	}

	// Bind JSON request body to booking struct
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update booking status to "cancelled" in the database
	result, bookingError := config.DB.Exec("UPDATE bookings SET status = 'cancelled' WHERE user_id = ? AND train_id = ?", booking.UserID, booking.TrainID)
	if bookingError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel booking"})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	// Check if any rows were affected
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	// Booking successfully cancelled
	c.JSON(http.StatusOK, gin.H{"message": "Booking cancelled successfully"})
}
