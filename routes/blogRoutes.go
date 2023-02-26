// Filename: blogRoutes.go
// This file contains the routes for the blog
package routes

import (
	"go-gorm-gauth/services"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func BlogRoutes(r *gin.Engine) {
	blog := r.Group("/blog")

	blog.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Blog routes",
		})
	})

	// Admin protected routes to create, update and delete posts
	blog.GET("/admin", services.AuthMiddleware(), func(c *gin.Context) {
		// Now we can access the user ID from the context (set in the AuthMiddleware)
		userID := c.MustGet("userID").(string)
		userIsAdmin := services.CheckIfUserIsAdmin(userID)

		if userIsAdmin {
			c.JSON(http.StatusOK, gin.H{
				"message": "Congrats you are an admin",
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}
	})

}
