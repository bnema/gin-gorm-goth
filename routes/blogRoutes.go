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

	// Public routes to view posts
	blog.GET("/", func(c *gin.Context) {
		// Get all posts from database
		posts := services.GetAllPosts()
		c.JSON(http.StatusOK, gin.H{
			"message": "Blog posts",
			"posts":   posts,
		})
	})

	// Unique route to view a single post by simplified title (replace espace with -)
	blog.GET("/:title", func(c *gin.Context) {
		title := c.Param("title")
		post := services.GetPostByTitle(title)
		c.JSON(http.StatusOK, gin.H{
			"message": "Blog post",
			"post":    post,
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
				"error": "Not enough power to see beyond this path",
			})
		}
	})
	// Create a new post
	blog.POST("/admin/createpost", services.AuthMiddleware(), func(c *gin.Context) {
		// Now we can access the user ID from the context (set in the AuthMiddleware)
		userID := c.MustGet("userID").(string)
		userIsAdmin := services.CheckIfUserIsAdmin(userID)

		if userIsAdmin {
			title := c.PostForm("title")
			content := c.PostForm("content")
			// Create a new post
			services.CreatePost(title, content, userID)
			c.JSON(http.StatusOK, gin.H{
				"message": "Post created",
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not enough power to see beyond this path",
			})
		}
	})
	// Update a post
	blog.POST("/admin/updatepost", services.AuthMiddleware(), func(c *gin.Context) {
		// Now we can access the user ID from the context (set in the AuthMiddleware)
		userID := c.MustGet("userID").(string)
		userIsAdmin := services.CheckIfUserIsAdmin(userID)
		id := c.PostForm("postID")

		if userIsAdmin {
			title := c.PostForm("title")
			content := c.PostForm("content")
			// Update a post
			services.UpdatePost(id, title, content, userID)
			c.JSON(http.StatusOK, gin.H{
				"message": "Post updated",
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not enough power to see beyond this path",
			})
		}
	})

	// Delete a post
	blog.POST("/admin/deletepost", services.AuthMiddleware(), func(c *gin.Context) {
		// Now we can access the user ID from the context (set in the AuthMiddleware)
		userID := c.MustGet("userID").(string)
		userIsAdmin := services.CheckIfUserIsAdmin(userID)
		id := c.PostForm("postID")

		if userIsAdmin {
			// Delete a post
			services.DeletePost(id, userID)
			c.JSON(http.StatusOK, gin.H{
				"message": "Post deleted",
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not enough power to see beyond this path",
			})
		}
	})
}
