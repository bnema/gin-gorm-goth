// Filename: authRoutes.go
// This file contains the routes for the authentication process
package routes

import (
	"go-gorm-gauth/services"
	"net/http"
	"os"
	"time"

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
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Error while authenticating",
			})
		}

		// Create new User and Account in the database
		_, err = services.CreateNewUser(authUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error creating new user",
			})
		}

		_, err = services.CreateNewAccount(authUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error creating new account",
			})
		}

		// Finally we create a session for the user and we return it in a Cookie HttpOnly
		session, err := services.CreateNewSession(authUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error creating new session",
			})
		}
		expirationDuration := session.Expires.Sub(time.Now())
		expirationSeconds := int(expirationDuration.Seconds())
		// Set cookie HTTPOnly with the session token / Expires now + 7 days (60*60*24*7)
		c.SetCookie("session_token", session.SessionToken, expirationSeconds, "/", (os.Getenv("DOMAIN")), false, true)
		c.SetCookie("session_id", session.ID, expirationSeconds, "/", (os.Getenv("DOMAIN")), false, true)

		// Return a CODE 200 and message "Authentication successful"
		c.JSON(http.StatusOK, gin.H{
			"message": "Authentication successful",
		})
	})
}
