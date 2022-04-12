package router

import (
	v1 "MedicalCare/api/v1"
	"MedicalCare/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "pong!",
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	apiUser := r.Group("/user")
	{
		apiUser.POST("/register", v1.UserRegister)
		apiUser.POST("/login", v1.UserLogin)
	}

	apiv1 := r.Group("/")
	{
		apiv1.GET("/user/info", v1.GetUserInfo)
	}

	// 404 信息返回
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "404 not found")
	})

	return r
}
