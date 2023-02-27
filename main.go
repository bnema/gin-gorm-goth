package main

import (
	"go-gorm-gauth/config"
	"go-gorm-gauth/routes"
	"runtime"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

var router = gin.Default()

func getRoutes() {
	routes.AuthRoutes(router)      // Auth routes with Goth and Gothic middleware for Oauth2 authentication
	routes.BlogRoutes(router)      // Blog routes for testing the database connection
	routes.BenchmarkRoutes(router) // Benchmark routes for testing the http server performance

	// handle errors if someone try to access a route that doesn't exist
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "404 - Not Found"})
	})

}

func main() {
	// Set the maximum number of CPUs that can be executing simultaneously and the maximum number of idle CPUs.
	runtime.GOMAXPROCS(runtime.NumCPU())
	config.InitDB() // Initialize database connection

	getRoutes()
	router.Run(":3000")

}
