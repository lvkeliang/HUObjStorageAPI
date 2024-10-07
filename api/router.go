package api

import (
	"HUObjStorageAPI/config"
	"HUObjStorageAPI/locate"
	"HUObjStorageAPI/objects"
	"HUObjStorageAPI/versions"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()

	// 定义路由和处理函数
	objectsGroup := r.Group("/objects")
	{
		objectsGroup.GET("/:name", objects.Get)
		objectsGroup.PUT("/:name", objects.Put)
		objectsGroup.DELETE("/:name", objects.Del)
	}

	r.GET("/locate/:name", locate.Handler)
	r.GET("/versions/:name", versions.Handler)

	// 启动服务
	//r.Run(os.Getenv("LISTEN_ADDRESS"))

	r.Run(config.Configs.ServerAddress)
}
