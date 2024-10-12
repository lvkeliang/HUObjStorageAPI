package temp

import (
	"HUObjStorageAPI/rs"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Head(c *gin.Context) {
	token := c.Param("token")
	stream, err := rs.NewRSResumablePutStreamFromToken(token)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusForbidden, gin.H{"info": "resumeToken invalid"})
		return
	}
	current := stream.CurrentSize()
	if current == -1 {
		c.JSON(http.StatusNotFound, gin.H{"info": "resume file not found"})
		return
	}
	//fmt.Println("HEAD: ", current)
	c.Header("content-length", fmt.Sprintf("%d", current))
	c.Status(http.StatusOK)
}
