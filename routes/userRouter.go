package routes

import (
	"go_event/controller"

	"github.com/gin-gonic/gin"
)

func UserRoutes(ginEngine *gin.Engine){
	ginEngine.GET("/users",controller.GetUsers())
	ginEngine.GET("/users/:user_id",controller.GetUser())
	ginEngine.POST("/users/signup",controller.SignUp())
	ginEngine.POST("/users/login",controller.Login())
	ginEngine.PATCH("/users/:user_id",controller.UpdateUser())
}