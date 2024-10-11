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

	status, err := storeObject(c.Request.Body, url.PathEscape(hash), size)
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

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	if locate.Exist(url.PathEscape(hash)) {
		return http.StatusOK, nil
	}

	stream, err := putStream(url.PathEscape(hash), size)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	reader := io.TeeReader(r, stream)

	// 从reader读数据时, 会同时读取r给reader并且写入到stream
	// 这里是实现传输给数据服务, 完成后计算hash
	d := util.CalculateHash(reader)

	if d != hash {
		// 哈希不对, 使用Commit把临时数据删除
		stream.Commit(false)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("object del failed")
		}
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch, calculated=%s, requested=%s", d, hash)
	}

	// 将临时数据转正
	stream.Commit(true)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("object commit failed")
	}
	return http.StatusOK, nil
}

func putStream(hash string, size int64) (*rs.RSPutStream, error) {
	servers := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	if len(servers) != rs.ALL_SHARDS {
		return nil, fmt.Errorf("cannot find enough dataServer")
	}

	return rs.NewRSPutStream(servers, hash, size)
}
