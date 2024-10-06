package api

import (
	"HUObjStorageAPI/config"
	"HUObjStorageAPI/locate"
	"HUObjStorageAPI/objects"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()

	// 定义路由和处理函数
	r.GET("/objects/:name", objects.Get)
	r.PUT("/objects/:name", objects.Put)
	r.GET("/locate/:name", locate.Handler)

	// 启动服务
	//r.Run(os.Getenv("LISTEN_ADDRESS"))

	r.Run(config.Configs.ServerAddress)
}
