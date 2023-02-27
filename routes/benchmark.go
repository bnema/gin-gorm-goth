// Filename benchmark.go
// This file is an attempt to benchmark the performance of the routes
package routes

import (
	"go-gorm-gauth/services"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func BenchmarkRoutes(r *gin.Engine) {
	benchmark := r.Group("/benchmark")

	// Benchmark route
	benchmark.GET("/", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "Benchmark route",
		})

		// Retrieve all posts by making a remote call to the database, which may increase the response time, but it is done to make the scenario more realistic.
		posts := services.GetAllPosts()
		c.JSON(200, gin.H{
			"posts": posts,
		})

	})

}
