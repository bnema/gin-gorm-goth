// Filename benchmark.go
// This file is an attempt to benchmark the performance of the routes
package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func BenchmarkRoutes(r *gin.Engine) {
	benchmark := r.Group("/benchmark")

	// Public routes to view posts
	benchmark.GET("/", func(c *gin.Context) {
		start := time.Now()
		// Votre logique de traitement ici
		elapsed := time.Since(start)
		c.JSON(http.StatusOK, gin.H{
			"response_time": elapsed.Milliseconds(),
		})
	})

}
