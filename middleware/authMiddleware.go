package middleware

import (
	"net/http"
	"restaraunt-management/helpers"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc{
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == ""{
			c.AbortWithStatus(401)
			return
		}
		claims, err := helpers.ValidateToken(clientToken)
		if err != ""{
			c.JSON(http.StatusForbidden, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid",claims.Uid)
		c.Next()
	}
}