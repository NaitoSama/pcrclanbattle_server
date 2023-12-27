package router

import (
	"github.com/gin-gonic/gin"
	"pcrclanbattle_server/common"
	"pcrclanbattle_server/config"
	"pcrclanbattle_server/server"
)

// router Add the request method and path here
func router(r *gin.Engine) {
	r.GET("/", func(context *gin.Context) {
		context.String(200, "test")
	})
	r.GET("/ws", server.Server.HandleConnection)
}

// RouterInit Http server startup
func RouterInit() {
	r := gin.New()
	router(r)
	common.Logln(0, "http server started")
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
