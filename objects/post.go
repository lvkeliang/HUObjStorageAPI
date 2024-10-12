package objects

import (
	"HUObjStorageAPI/es"
	"HUObjStorageAPI/heartbeat"
	"HUObjStorageAPI/locate"
	"HUObjStorageAPI/rs"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func Post(c *gin.Context) {
	objectName := c.Param("name")
	size, err := strconv.ParseInt(c.GetHeader("size"), 0, 64)

	if size <= 0 || err != nil {
		log.Println("size header invalid")
		c.JSON(http.StatusBadRequest, gin.H{"info": "size header invalid"})
		return
	}

	hash := GetHash(c.GetHeader("digest"))

	if hash == "" {
		log.Println("missing object hash in digest header")
		c.JSON(http.StatusBadRequest, gin.H{"info": "missing object hash in digest header"})
		return
	}

	if locate.Exist(url.PathEscape(hash)) {
		err = es.AddVersion(objectName, hash, size)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"info": "add version failed"})
		} else {
			c.JSON(http.StatusOK, gin.H{"info": "success"})
		}
		return
	}

	ds := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	if len(ds) != rs.ALL_SHARDS {
		log.Println("cannot find enough dataServer")
		c.JSON(http.StatusServiceUnavailable, gin.H{"info": "dataServer not enough"})
		return
	}

	stream, err := rs.NewRSResumablePutStream(ds, objectName, url.PathEscape(hash), size)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"info": "get put stream failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"info": "success", "location": "/temp/" + url.PathEscape(stream.ToToken())})
}
