// Filename: authRoutes.go
// This file contains the routes for the authentication process
package routes

import (
	"fmt"
	"go-gorm-gauth/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")

	auth.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Auth routes",
		})
	})

	// Initialize authentication process
	auth.GET("/login", func(c *gin.Context) {
		// Set providers
		goth.UseProviders(
			discord.New(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"), os.Getenv("AUTH_CALLBACK_URL")),

			// Other providers here...
		)
		// Start authentication process
		gothic.BeginAuthHandler(c.Writer, c.Request)

	})

	auth.GET("/callback", func(c *gin.Context) {
		// Complete authentication process
		authUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Error while authenticating",
			})
		}

		services.CreateNewUser(authUser)

		// Return user data as JSON
		c.JSON(http.StatusOK, authUser)
	})
}
