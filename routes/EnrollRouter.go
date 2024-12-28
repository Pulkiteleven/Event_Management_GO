package routes

import (
	"go_event/controller"

	"github.com/gin-gonic/gin"
)

func EnrollRoutes(GinEngine *gin.Engine){
	GinEngine.POST("enroll",controller.EnrollUser())
	GinEngine.PATCH("enroll/:user_id/:event_id",controller.ApprovedUser())
	GinEngine.GET("enroll/:event_id",controller.GetEnrollUser())
}