package objects

import (
	"HUObjStorageAPI/es"
	"HUObjStorageAPI/locate"
	"HUObjStorageAPI/objectstream"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func Get(c *gin.Context) {
	objectName := c.Param("name")
	versionID := c.Query("version")

	version := 0
	var err error
	if len(versionID) != 0 {
		version, err = strconv.Atoi(versionID)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"info": "version invalid"})
			return
		}
	}

	meta, err := es.GetMetadata(objectName, version)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"info": "find resource failed"})
		return
	}

	if meta.Hash == "" {
		c.JSON(http.StatusNotFound, gin.H{"info": "resource not found"})
		return
	}

	object := url.PathEscape(meta.Hash)

	stream, err := GetStream(object)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"info": "resource not found"})
	}
	io.Copy(c.Writer, stream)
}

func GetStream(object string) (io.Reader, error) {
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate failed", object)
	}
	return objectstream.NewGetStream(server, object)
}
