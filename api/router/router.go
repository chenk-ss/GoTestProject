package router

import (
	"goTestProject/api/handler"
	"goTestProject/tools"
	"net/http"

	_ "goTestProject/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Register() *gin.Engine {
	r := gin.Default()
	r.Use(CorsMiddleware())
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	initUserRouter(r)
	initFileRouter(r)
	go handler.ChannelConsume()
	r.NoRoute(func(c *gin.Context) {
		tools.FailWithMsg(c, "please check request url !")
	})
	return r
}

func initUserRouter(r *gin.Engine) {
	userGroup := r.Group("/user")
	userGroup.GET("/info", handler.UserInfo)
	userGroup.POST("/register", handler.Register)
	//checkAuthInfo
	// userGroup.Use(CheckSessionId())
	{
		// userGroup.POST("/checkAuth", handler.CheckAuth)
	}
}

func initFileRouter(r *gin.Engine) {
	userGroup := r.Group("/file")
	userGroup.GET("/collect", handler.FileCollect)
	userGroup.GET("/random", handler.Random)
	userGroup.GET("/channel", handler.ChannelTest)
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		var openCorsFlag = true
		if openCorsFlag {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, nil)
		}
		c.Next()
	}
}
