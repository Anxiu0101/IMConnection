package router

import (
	v1 "IMConnection/api/v1"
	"IMConnection/docs"
	"IMConnection/middleware"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"

	// util router
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "pong!",
		})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// no token router
	apiUser := r.Group("/user")
	{
		apiUser.POST("/register", v1.UserRegister)
		apiUser.POST("/login", v1.UserLogin)
		apiUser.GET("/token", v1.RefreshAccessToken)
	}

	// token router
	apiv1 := r.Group("/")
	apiv1.Use(middleware.JWT())
	{
		apiv1.POST("/user/password", v1.ResetUserPassword)
		apiv1.GET("/user/info", v1.GetUserInfo)
		apiv1.POST("/user/info", v1.UpdateUserInfo)

		apiv1.POST("/group", v1.CreateGroup)
		apiv1.PUT("/group/member", v1.InviteMember)

		apiv1.GET("/chat/:receiver", v1.Chat)
	}

	// 404 信息返回
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "404 not found")
	})

	return r
}
