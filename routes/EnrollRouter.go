package routes

import (
	"go_event/controller"

	"github.com/gin-gonic/gin"
)

func EnrollRoutes(GinEngine *gin.Engine){
	GinEngine.POST("Enroll",controller.EnrollUser())
	GinEngine.PATCH("Enroll/:user_id/:event_id",controller.ApprovedUser())
	GinEngine.GET("Enroll/:event_id",controller.GetEnrollUser())
}