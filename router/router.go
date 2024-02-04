package router

import (
	"github.com/gin-gonic/gin"
	"pcrclanbattle_server/common"
	"pcrclanbattle_server/config"
	"pcrclanbattle_server/server"
)

// router Add the request method and path here
func router(r *gin.Engine) {
	r.Static("/pic", "./pic")
	main := r.Group("/")
	main.Use(common.RandTokenSet)
	{
		v1 := main.Group("v1")
		v1.Use(common.JWTAuthentication)
		{
			v1.GET("/ws", server.Server.HandleConnection)
			v1.GET("/records", server.GetRecords)
			v1.POST("/uploadbosspic", server.UploadBossPic)
		}
		main.POST("login", server.Login)
		main.POST("register", server.Register)
		main.POST("userinfo", server.GetUserInfoFromJWT)
		main.GET("test", func(context *gin.Context) {
			context.String(200, "test")
		})
	}
}

// RouterInit Http server startup
func RouterInit() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	router(r)
	common.Logln(0, "http server started")
	println("Server started on \"IPv4 OR IPv6:" + config.Config.General.HttpPort + "\"!")
	err := r.Run(":" + config.Config.General.HttpPort)
	if err != nil {
		common.ErrorHandle(err)
		return
	}
	defer func() {
		err1 := recover()
		if err1 != nil {
			common.Logln(2, err1)
		}
	}()

}
