package main

import (
	"go-gorm-gauth/config"
	"go-gorm-gauth/routes"

	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func getRoutes() {
	routes.AuthRoutes(router) // Auth routes with Goth and Gothic middleware for Oauth2 authentication
	routes.BlogRoutes(router) // Blog routes
	routes.BenchmarkRoutes(router)
}

func main() {

	config.InitDB() // Initialize database connection

	getRoutes()
	router.Run(":3000")

}
