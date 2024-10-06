package objects

import (
	"HUObjStorageAPI/heartbeat"
	"HUObjStorageAPI/objectstream"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func Put(c *gin.Context) {
	objectName := c.Param("name")
	status, err := storeObject(c.Request.Body, objectName)
	if err != nil {
		log.Println(err)
	}
	if status != http.StatusOK {
		c.JSON(status, gin.H{"info": err.Error()})
	} else {
		c.JSON(status, gin.H{"info": "success"})
	}

}

func storeObject(r io.Reader, objectName string) (int, error) {
	stream, err := putStream(objectName)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}

	io.Copy(stream, r)
	err = stream.Close()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func putStream(objectName string) (*objectstream.PutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}

	return objectstream.NewPutStream(server, objectName), nil
}
