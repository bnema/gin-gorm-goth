// Filename: authRoutes.go
// This file contains the routes for the authentication process
package routes

import (
	"context"
	"go-gorm-gauth/services"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/google"
)

// Create a slice of provider names
var providerNames []string

func AuthRoutes(r *gin.Engine) {
	goth.UseProviders(
		discord.New(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"), os.Getenv("AUTH_CALLBACK_URL")),

		// Add more providers here
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), os.Getenv("AUTH_CALLBACK_URL")),
	)
	for _, provider := range goth.GetProviders() {
		// Append the provider's name to the slice
		providerNames = append(providerNames, provider.Name())
	}

	// Auth routes
	auth := r.Group("/auth")
	auth.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Auth routes",
		})
	})

	// Initialize authentication process
	auth.GET("/login", func(c *gin.Context) {

		// json all the providers available
		c.JSON(http.StatusOK, gin.H{
			"message": "Available providers",
			// For each provider, we return the name and the url
			"providers": providerNames,
		})

	})

	// Dynamic route with the provider name
	auth.GET("/login/:provider", func(c *gin.Context) {
		// Get the provider from the url
		provider := c.Param("provider")

		// Check if the provider exists in the slice
		for _, name := range providerNames {
			if name == provider {
				// If the provider exists, we start the authentication process
				gothic.BeginAuthHandler(c.Writer, c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", provider)))
				return
			}
		}
		// If the provider doesn't exist, we return an error
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Provider not found",
		})
	})

	// Logout route (delete session from database and delete cookie)
	auth.GET("/logout", func(c *gin.Context) {
		// Get the session token from the cookie
		sessionToken, err := c.Cookie("session_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Error getting session token from cookie",
			})
		}
		// Delete the session from the database
		err = services.DeleteSession(sessionToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error deleting session from database",
			})
		}
		// Delete the cookie
		c.SetCookie("session_token", "", -1, "/", (os.Getenv("DOMAIN")), true, true)
		c.SetCookie("session_id", "", -1, "/", (os.Getenv("DOMAIN")), true, true)
		c.JSON(http.StatusOK, gin.H{
			"message": "Logged out",
		})
	})

	auth.GET("/callback", func(c *gin.Context) {
		// Complete authentication process
		authUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Error while authenticating",
			})
		}
		// Check if the user already exist in the database
		userExist := services.CheckIfUserExists(authUser.Email)

		// If the user already exists, we update it
		if userExist {
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error updating user",
				})
			}
			// Get the account from the database
			providerAccountID := authUser.UserID
			account, err := services.GetAccountByProviderAccountID(providerAccountID)
			// If the account already exists, we update it
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error getting account",
				})
			} else {
				_, err = services.UpdateAccount(account)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "Error updating account",
					})
				}
			}

			// If the user already exists, we create a new session for him and we return it in a Cookie HttpOnly
			session, err := services.CreateNewSession(authUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error creating new session",
				})
			}
			expirationDuration := session.Expires.Sub(time.Now())
			expirationSeconds := int(expirationDuration.Seconds())
			// Set cookie HTTPOnly with the session token / Expires now + 7 days (60*60*24*7)
			c.SetCookie("session_token", session.SessionToken, expirationSeconds, "/", (os.Getenv("DOMAIN")), true, true)
			c.SetCookie("session_id", session.ID, expirationSeconds, "/", (os.Getenv("DOMAIN")), false, true)

			// Return a CODE 200 and message "Authentication successful"
			c.JSON(http.StatusOK, gin.H{
				"message": "Authentication successful",
			})
			return
		}

		// If the user doesn't exist, we create a new user and a new account in the database
		if !userExist {

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

		}
	})
}
