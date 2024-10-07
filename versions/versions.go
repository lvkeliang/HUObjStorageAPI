package versions

import (
	"HUObjStorageAPI/es"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Handler(c *gin.Context) {
	objectName := c.Param("name")

	from := 0
	size := 1000

	for {
		metas, err := es.SearchAllVersions(objectName, from, size)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"info": err.Error()})
			return
		}

		for i := range metas {
			b, _ := json.Marshal(metas[i])
			c.Writer.Write(b)
			c.Writer.Write([]byte("\n"))
		}
		if len(metas) != size {
			return
		}
		from += size
	}
}
