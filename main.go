package main

import (
	"go-gorm-gauth/config"
	"go-gorm-gauth/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

var router = gin.Default()

func getRoutes() {
	routes.AuthRoutes(router) // Auth routes with Goth and Gothic middleware for Oauth2 authentication
	routes.BlogRoutes(router) // Blog routes
	routes.BenchmarkRoutes(router)
}

func main() {
	// Set the maximum number of CPUs that can be executing simultaneously and the maximum number of idle CPUs.
	// runtime.GOMAXPROCS(runtime.NumCPU())
	config.InitDB() // Initialize database connection

	getRoutes()
	router.Run(":3000")

}
