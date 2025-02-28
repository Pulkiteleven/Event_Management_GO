package middleware

import (
	"fmt"
	"go_event/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc{
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == ""{
			c.JSON(http.StatusInternalServerError,gin.H{"error":fmt.Sprintf("No Authorization Provided")})
			c.Abort()
			return
		}

		claims,err := helpers.ValidateToken(clientToken)
		if err != ""{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err})
			c.Abort()
			return
		}

		c.Set("email",claims.Email)
		c.Set("user",claims.User_name)
		c.Set("uid",claims.Uid)

		c.Next()
	}
}