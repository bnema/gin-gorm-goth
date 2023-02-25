package main

import (
	"go-gorm-gauth/routes"

	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func getRoutes() {
	routes.AuthRoutes(router) // Auth routes with Goth and Gothic middleware for Oauth2 authentication
}

func main() {
	getRoutes()
	router.Run(":3000")

}
