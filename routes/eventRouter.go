package routes

import (
	"go_event/controller"

	"github.com/gin-gonic/gin"
)

func EventRoutes(ginEngine *gin.Engine){
	ginEngine.POST("/events",controller.CreateEvent())
	ginEngine.GET("/events",controller.GetEvents())
	ginEngine.GET("/events/:user_id",controller.GetUserEnrolledEvents())
}

