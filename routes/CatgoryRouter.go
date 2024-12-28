package routes

import (
	"go_event/controller"

	"github.com/gin-gonic/gin"
)

func CategoryRoutes(GinEngine *gin.Engine){
	GinEngine.POST("category",controller.CreateCategroy())
	GinEngine.GET("category",controller.GetCategory())
}