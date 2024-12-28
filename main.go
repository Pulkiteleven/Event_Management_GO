package main

import (
	"go_event/middleware"
	"go_event/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == ""{
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())


	routes.UserRoutes(router)

	router.Use(middleware.Authentication())

	routes.EventRoutes(router)
	routes.EnrollRoutes(router)
	routes.AttendanceRoutes(router)
	routes.CategoryRoutes(router)
	router.Run(":" + port)

}