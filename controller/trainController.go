package controllers

import (
	"SanjaySinghRajpoot/ticketBooking/config"
	"SanjaySinghRajpoot/ticketBooking/models"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CheckAvailableTrains(c *gin.Context) {

	var checkTrains struct {
		To   string `json:"to" binding:"required"`
		From string `json:"from" binding:"required"`
	}

	if err := c.ShouldBindJSON(&checkTrains); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	rows, err := config.DB.Query("SELECT * FROM trains WHERE \"to\"=$1 AND \"from\"=$2 LIMIT 100", checkTrains.To, checkTrains.From)
	if err != nil {
		fmt.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	defer rows.Close()

	var allTains []models.Train

	for rows.Next() {
		var train models.Train
		if err := rows.Scan(&train.ID, &train.Name, &train.DepartureTime, &train.ArrivalTime, &train.From, &train.To, &train.TotalSeats, &train.Fare, &train.CreatedAt); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}

		allTains = append(allTains, train)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"trains": allTains,
	})
}

func CheckAvailableSeats(c *gin.Context) {
	trainID := c.Query("train_id")
	if trainID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Train ID is required",
		})
		return
	}

	var availableSeats int
	err := config.DB.QueryRow("SELECT total_seats FROM trains WHERE id=$1", trainID).Scan(&availableSeats)
	if err != nil {
		fmt.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check available seats",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"available_seats": availableSeats,
	})
}

func SaveTrain(c *gin.Context) {

	var newTrain struct {
		Name          string    `json:"name" binding:"required"`
		DepartureTime time.Time `json:"departure_time" binding:"required"`
		ArrivalTime   time.Time `json:"arrival_time" binding:"required"`
		From          string    `json:"from" binding:"required"`
		To            string    `json:"to" binding:"required"`
		TotalSeats    int       `json:"total_seats" binding:"required"`
		Fare          int       `json:"fare" binding:"required"`
		CreatedAt     time.Time `json:"created_at"`
		AdminKey      string    `json:"admin_key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&newTrain); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check admin api key
	adminKey := os.Getenv("ADMIN_KEY")

	if newTrain.AdminKey != adminKey {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Not authorized to perform this action"})
		return
	}

	// Set the created_at field to the current time
	newTrain.CreatedAt = time.Now()

	// Execute the SQL INSERT statement
	_, err := config.DB.Exec("INSERT INTO trains (name, departure_time, arrival_time, \"from\", \"to\", total_seats, fare, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		newTrain.Name, newTrain.DepartureTime, newTrain.ArrivalTime, newTrain.From, newTrain.To, newTrain.TotalSeats, newTrain.Fare, newTrain.CreatedAt)
	if err != nil {
		fmt.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add train"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Train added successfully"})

}

func DeleteTrain(c *gin.Context) {
	// Extract train ID from query parameter
	trainID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid train ID"})
		return
	}

	adminKeyInput := c.Query("admin_key")

	// check admin api key
	adminKey := os.Getenv("ADMIN_KEY")

	if adminKeyInput != adminKey {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Not authorized to perform this action"})
		return
	}

	// Update train record with deleted_at timestamp
	result, err := config.DB.Exec("UPDATE trains SET deleted_at = $1 WHERE id = $2", time.Now(), trainID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete train"})
		return
	}

	delrow, _ := result.RowsAffected()

	// Check if any rows were affected
	if delrow == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Train not found"})
		return
	}

	// Train successfully deleted
	c.JSON(http.StatusOK, gin.H{"message": "Train deleted successfully"})
}
