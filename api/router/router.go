package router

import (
	"bytes"
	"goTestProject/api/handler"
	"goTestProject/tools"
	"log"
	"net/http"
	"time"

	_ "goTestProject/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Register() *gin.Engine {
	// Logging to a file.
	// f, _ := os.Create("gin.log")
	// gin.DefaultWriter = io.MultiWriter(f)

	r := gin.Default()
	r.Use(CorsMiddleware())
	// r.Use(ginBodyLogMiddleware)
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	initUserRouter(r)
	initFileRouter(r)
	initKafkaRouter(r)
	go handler.ChannelConsume()
	r.NoRoute(func(c *gin.Context) {
		tools.FailWithMsg(c, "please check request url !")
	})
	return r
}

func initUserRouter(r *gin.Engine) {
	userGroup := r.Group("/user")
	userGroup.POST("/register", handler.Register)
	userGroup.POST("/login", handler.Login)
	//checkAuthInfo
	userGroup.Use(CheckSessionId())
	{
		userGroup.GET("/info", handler.UserInfo)
		// userGroup.POST("/checkAuth", handler.CheckAuth)
	}
}

func initFileRouter(r *gin.Engine) {
	userGroup := r.Group("/file")
	userGroup.GET("/collect", handler.FileCollect)
	userGroup.GET("/random", handler.Random)
	userGroup.GET("/channel", handler.ChannelTest)
	userGroup.POST("/upload", handler.UploadFile)
}

func initKafkaRouter(r *gin.Engine) {
	kafkaGroup := r.Group("/kafka")
	kafkaGroup.POST("/push", handler.PushMsg)
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

type FormCheckSessionId struct {
	AuthToken string `form:"authToken" json:"authToken" binding:"required"`
}

func CheckSessionId() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formCheckSessionId FormCheckSessionId
		if err := c.ShouldBindHeader(&formCheckSessionId); err != nil {
			c.Abort()
			tools.ResponseWithCode(c, tools.CodeSessionError, nil, nil)
			return
		}
		authToken := formCheckSessionId.AuthToken
		if !handler.CheckToken(authToken) {
			c.Abort()
			tools.ResponseWithCode(c, tools.CodeSessionError, nil, nil)
			return
		}
		c.Next()
		return
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// gin request log
func ginBodyLogMiddleware(c *gin.Context) {
	t := time.Now()
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()
	latency := time.Since(t)
	log.Printf(`
		Lantency: %s
		Request url: %s
		Request body: %s
		Response body: %s`,
		latency.String(), c.Request.RequestURI, c.Request.Body, blw.body.String())
}
