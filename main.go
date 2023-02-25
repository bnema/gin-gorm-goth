package main

import (
	"go-gorm-gauth/routes"

	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func getRoutes() {
	// AuthRoutes(router)
	routes.AuthRoutes(router)

}

func main() {

	getRoutes()
	router.Run(":3000")

}
