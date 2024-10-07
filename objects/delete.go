package objects

import (
	"HUObjStorageAPI/es"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Del(c *gin.Context) {
	objectName := c.Param("name")
	version, err := es.SearchLatestVersion(objectName)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"info": err.Error()})
		return
	}
	err = es.PutMetadata(objectName, version.Version+1, 0, "")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"info": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"info": "success"})
}
