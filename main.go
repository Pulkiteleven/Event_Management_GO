package main

import (
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
	routes.EventRoutes(router)
	routes.EnrollRoutes(router)
	routes.AttendanceRoutes(router)

	router.Run(":" + port)

}