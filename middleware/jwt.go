package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"IMConnection/pkg/e"
	"IMConnection/pkg/util"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int

		code = e.Success
		token := c.GetHeader("Authorization")
		if token == "" {
			code = e.InvalidParams
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.Error
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.Error
			}
		}

		if code != e.Success {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
