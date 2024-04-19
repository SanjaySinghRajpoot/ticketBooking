package middleware

import (
	"SanjaySinghRajpoot/ticketBooking/config"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthUser struct {
	ID    uint   `json:"ID"`
	Name  string `json:"Name"`
	Email string `json:"Email"`
}

func RequireAuth(c *gin.Context) {
	// Get the cookie from the request
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// Decode and validate it
	// Parse and takes the token string and a function for look
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the expiration time
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Find the user with token sub
		var userId int
		var username string
		var email string
		// initializers.DB.Find(&user, claims["sub"])

		rows, err := config.DB.Query("SELECT id, name, email FROM users WHERE id=$1", claims["sub"])
		if err != nil {
			fmt.Println("Error getting records:", err)
			return
		}
		defer rows.Close()

		if rows.Next() {
			if err := rows.Scan(&userId, &username, &email); err != nil {
				fmt.Println("Error scanning row:", err)
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		authUser := AuthUser{
			ID:    uint(userId),
			Name:  username,
			Email: email,
		}

		// Attach the user to request
		c.Set("authUser", authUser)

		// Continue
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
