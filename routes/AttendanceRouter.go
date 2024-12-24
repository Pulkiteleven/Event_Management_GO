package routes

import (
	"go_event/controller"

	"github.com/gin-gonic/gin"
)

func AttendanceRoutes(GinEngine *gin.Engine){
	GinEngine.POST("/attendance",controller.MarkAttendance())
	GinEngine.GET("/attendance/:event_id",controller.GetAttendance())
}