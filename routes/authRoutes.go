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
		userFromDB, err := services.GetUserByEmail(authUser.Email)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error checking if user exists",
			})
		} else if userFromDB != nil {
			// Update the user and the account in the database
			_, err = services.UpdateUser(userFromDB)

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
		if userFromDB == nil {

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
