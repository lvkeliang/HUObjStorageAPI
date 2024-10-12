package objects

import (
	"HUObjStorageAPI/es"
	"HUObjStorageAPI/heartbeat"
	"HUObjStorageAPI/locate"
	"HUObjStorageAPI/rs"
	"HUObjStorageAPI/util"
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
	offset := util.GetOffset(c.GetHeader("range"))

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

	objectHash := url.PathEscape(meta.Hash)

	stream, err := GetStream(objectHash, meta.Size)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"info": "resource not found"})
		return
	}

	if offset != 0 {
		stream.Seek(offset, io.SeekCurrent)
		c.JSON(http.StatusOK, gin.H{"info": "success", "content-range": fmt.Sprintf("bytes %d-%d/%d", offset, meta.Size-1, meta.Size)})
	}

	_, err = io.Copy(c.Writer, stream)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"info": "resource not found"})
		return
	}
	stream.Close()
}

func GetStream(objectHash string, size int64) (*rs.RSGetStream, error) {
	locateInfo := locate.Locate(objectHash)
	if len(locateInfo) < rs.DATA_SHARDS {
		return nil, fmt.Errorf("object %s locate failed, result %v", objectHash, locateInfo)
	}

	dataServers := make([]string, 0)
	if len(locateInfo) != rs.ALL_SHARDS {
		dataServers = heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS-len(locateInfo), locateInfo)
	}
	return rs.NewRSGetStream(locateInfo, dataServers, objectHash, size)
}
