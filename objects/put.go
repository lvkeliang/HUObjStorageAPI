package objects

import (
	"HUObjStorageAPI/es"
	"HUObjStorageAPI/heartbeat"
	"HUObjStorageAPI/objectstream"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func GetHash(digest string) string {
	if len(digest) < 9 {
		return ""
	}
	if digest[:8] != "SHA-256=" {
		return ""
	}
	return digest[8:]
}

func Put(c *gin.Context) {
	objectName := c.Param("name")
	hash := GetHash(c.GetHeader("digest"))
	size, err := strconv.ParseInt(c.GetHeader("content-length"), 0, 64)

	//fmt.Println(c.GetHeader("digest"))
	//fmt.Println(c.GetHeader("content-length"))

	if size <= 0 || err != nil {
		log.Println("size header invalid")
		c.JSON(http.StatusBadRequest, gin.H{"info": "missing object hash in digest header"})
		return
	}

	if hash == "" {
		log.Println("missing object hash in digest header")
		c.JSON(http.StatusBadRequest, gin.H{"info": "missing object hash in digest header"})
		return
	}

	status, err := storeObject(c.Request.Body, url.PathEscape(hash))
	if err != nil {
		log.Println(err)
		c.JSON(status, gin.H{"info": "store object failed"})
		return
	}
	if status != http.StatusOK {
		c.JSON(status, gin.H{"info": "store object failed"})
		return
	}

	err = es.AddVersion(objectName, hash, size)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"info": "store object failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"info": "success"})
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
