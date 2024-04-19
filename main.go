package main

import (
	"SanjaySinghRajpoot/ticketBooking/config"
	"SanjaySinghRajpoot/ticketBooking/router"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	config.ConnectDB()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	fmt.Println("Hello auth")
	r := gin.Default()
	router.GetRoute(r)

	r.Run()
}
