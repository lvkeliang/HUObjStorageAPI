package objects

import (
	"HUObjStorageAPI/locate"
	"HUObjStorageAPI/objectstream"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func Get(c *gin.Context) {
	objectName := c.Param("name")
	stream, err := GetStream(objectName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"info": "resource not found"})
	}
	io.Copy(c.Writer, stream)
}

func GetStream(objectName string) (io.Reader, error) {
	server := locate.Locate(objectName)
	if server == "" {
		return nil, fmt.Errorf("object %s locate failed", objectName)
	}
	return objectstream.NewGetStream(server, objectName)
}
