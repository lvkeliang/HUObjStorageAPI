package locate

import (
	"HUObjStorageAPI/config"
	"HUObjStorageAPI/rabbitmq"
	"HUObjStorageAPI/rs"
	"HUObjStorageAPI/types"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Handler(c *gin.Context) {

	info := Locate(c.Param("name"))

	if len(info) == 0 {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}

	b, err := json.Marshal(info)
	if err != nil {
		// 如果 JSON 序列化失败，返回 500 错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize file info"})
		return
	}

	// 返回文件信息的 JSON 响应
	c.Data(http.StatusOK, "application/json", b)
}

func Locate(name string) (locateInfo map[int]string) {
	q := rabbitmq.New(config.Configs.Rabbitmq.RabbitmqServer)
	q.Publish("dataServers", name)

	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()

	locateInfo = make(map[int]string)
	for i := 0; i < rs.ALL_SHARDS; i++ {
		msg := <-c
		if len(msg.Body) == 0 {
			return
		}
		var info types.LocateMessage
		json.Unmarshal(msg.Body, &info)
		locateInfo[info.Id] = info.Addr
	}
	return
}

func Exist(name string) bool {
	return len(Locate(name)) >= rs.DATA_SHARDS
}
