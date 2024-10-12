package api

import (
	"HUObjStorageAPI/config"
	"HUObjStorageAPI/locate"
	"HUObjStorageAPI/objects"
	"HUObjStorageAPI/temp"
	"HUObjStorageAPI/versions"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()

	// 定义路由和处理函数
	objectsGroup := r.Group("/objects")
	{
		// 上传对象，若对象已存在则更新版本信息
		// 使用流程：客户端需要上传一个对象，首先计算对象的hash并通过POST接口创建临时文件，之后用PUT接口上传对象数据。
		objectsGroup.POST("/:name", objects.Post)

		// 获取对象
		// 使用流程：客户端使用GET请求获取指定name的对象，可以通过可选的查询参数获取特定版本，支持Range Header进行断点续传。
		objectsGroup.GET("/:name", objects.Get)

		// 上传对象内容
		// 使用流程：客户端上传数据时使用PUT请求，传递文件内容、大小和hash信息；适用于初次上传和覆盖现有对象。
		objectsGroup.PUT("/:name", objects.Put)

		// 删除对象
		// 使用流程：客户端发送DELETE请求来删除指定对象，会创建一个新版本且内容为空，表示逻辑删除。
		objectsGroup.DELETE("/:name", objects.Del)
	}

	tempGroup := r.Group("/temp")
	{
		// 上传对象的分片数据，用于断点续传
		// 使用流程：当对象较大时，通过POST接口获取token，然后使用该token调用PUT接口上传数据，支持部分数据上传。
		tempGroup.PUT("/:token", temp.Put)

		// 获取已上传的文件片段大小，用于断点续传时确定上传位置
		// 使用流程：客户端使用HEAD请求检查已上传的文件部分，以确定从哪里继续上传。
		tempGroup.HEAD("/:token", temp.Head)
	}

	// 查询对象的位置
	// 使用流程：客户端通过hash值查找对象位置，主要用于判断对象是否已存在于系统中。
	r.GET("/locate/:hash", locate.Handler)

	// 查询对象的所有版本
	// 使用流程：客户端可以通过GET请求获取指定name的所有版本的元数据，用于版本控制和管理。
	r.GET("/versions/:name", versions.Handler)

	// 启动服务
	// 配置监听地址，启动Gin HTTP服务
	r.Run(config.Configs.ServerAddress)
}
