package services

import (
	"fmt"
	"go-gorm-gauth/config"
	"go-gorm-gauth/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// Create a AuthMiddleware function that will be used to protect routes (with Gin)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the session token from the request
		sessionToken := c.GetHeader("Authorization")
		// If the session token is empty, return an error
		if sessionToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		// Export the claims from the JWT token
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(sessionToken, claims, func(token *jwt.Token) (interface{}, error) {

			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			fmt.Println("Error parsing JWT token:", err)
			// If the JWT token is invalid, return an error
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session token"})
			c.Abort()
			return
		}

		// Get the user ID from the claims
		userID := claims["id"].(string)
		// Get the user from the database
		user := &models.User{}
		err = config.DB.Where("id = ?", userID).First(user).Error
		if err != nil {
			fmt.Println("Error getting user by id:", err)
			// If no user is found with the specified ID, return an error
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Unauthorized",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			c.Abort()
			return
		}
		// Set the user ID in the context
		c.Set("userID", userID)
		c.Next()
	}
}
