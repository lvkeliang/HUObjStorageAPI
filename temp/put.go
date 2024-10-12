package temp

import (
	"HUObjStorageAPI/es"
	"HUObjStorageAPI/locate"
	"HUObjStorageAPI/rs"
	"HUObjStorageAPI/util"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
)

func Put(c *gin.Context) {
	token := c.Param("token")
	stream, err := rs.NewRSResumablePutStreamFromToken(token)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusForbidden, gin.H{"info": "get stream failed"})
		return
	}

	current := stream.CurrentSize()
	if current == -1 {
		c.JSON(http.StatusNotFound, gin.H{"info": "tempFile not found"})
		return
	}

	offset := util.GetOffset(c.GetHeader("range"))
	if current != offset {
		c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"info": "range header invalid"})
		return
	}

	bytes := make([]byte, rs.BLOCK_SIZE)
	for {
		n, err := io.ReadFull(c.Request.Body, bytes)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"info": "upload failed"})
			return
		}
		current += int64(n)
		if current > stream.Size {
			stream.Commit(false)
			log.Println("resumable put exceed size")
			c.JSON(http.StatusForbidden, gin.H{"info": "resumable put exceed size"})
			return
		}

		if n != rs.BLOCK_SIZE && current != stream.Size {
			return
		}
		stream.Write(bytes[:n])
		if current == stream.Size {
			stream.Flush()
			getStream, err := rs.NewRSResumableGetStream(stream.Servers, stream.Uuids, stream.Size)
			hash := util.CalculateHash(getStream)
			if hash != stream.Hash {
				stream.Commit(false)
				log.Println("resumable put done but hash mismatch")
				c.JSON(http.StatusForbidden, gin.H{"info": "resumable put done but hash mismatch"})
				return
			}
			if locate.Exist(url.PathEscape(hash)) {
				stream.Commit(false)
			} else {
				stream.Commit(true)
			}
			err = es.AddVersion(stream.Name, stream.Hash, stream.Size)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"info": "metadata add failed"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"info": "success"})
			return
		}
	}
}
