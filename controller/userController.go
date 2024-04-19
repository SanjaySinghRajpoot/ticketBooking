package controllers

import (
	"SanjaySinghRajpoot/ticketBooking/config"
	"SanjaySinghRajpoot/ticketBooking/models"
	"SanjaySinghRajpoot/ticketBooking/validations"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Signup function is used to create a user or signup a user
func Signup(c *gin.Context) {
	// Get the name, email and password from request
	var userInput struct {
		Name     string `json:"name" binding:"required,min=2,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"err": errs,
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Email unique validation
	if validations.IsUniqueValue("users", "email", userInput.Email) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"validations": map[string]interface{}{
				"Email": "The email is already exist!",
			},
		})
		return
	}

	// Hash the password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	var userID int64

	insertErr := config.DB.QueryRow("INSERT INTO users (name, email, password, created_at) VALUES ($1, $2, $3, $4) RETURNING id", userInput.Name, userInput.Email, userInput.Password, time.Now()).Scan(&userID)
	if insertErr != nil {
		fmt.Println("Error inserting records:", insertErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr})
		return
	}

	user := models.User{
		ID:       userID,
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: string(hashPassword),
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// Login function is used to log in a user
func Login(c *gin.Context) {
	// Get the email and password from the request
	var userInput struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if c.ShouldBindJSON(&userInput) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Find the user by email
	// var user models.User
	// initializers.DB.First(&user, "email = ?", userInput.Email)

	rows, err := config.DB.Query("SELECT password, id FROM users WHERE email=$1", userInput.Email)
	if err != nil {
		fmt.Println("Error getting records:", err)
		return
	}

	defer rows.Close()

	var password string
	var userID string

	if rows.Next() {
		if err := rows.Scan(&password, &userID); err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
	}

	// Compare the password with user hashed password
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(userInput.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign in and get the complete encoded token as a string using the .env secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// Set expiry time and send the token back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status": "login succesful",
	})
}

func Logout(c *gin.Context) {
	// Clear the cookie
	c.SetCookie("Authorization", "", 0, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"successMessage": "Logout successful",
	})
}
