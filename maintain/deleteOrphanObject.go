package maintain

import (
	"HUObjStorageAPI/config"
	"HUObjStorageAPI/es"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// DeleteOrphanObject 删除没有元数据引用的对象数据
func DeleteOrphanObject() {
	files, _ := filepath.Glob(config.Configs.StorageRoot + "/objects/*")

	for i := range files {
		hash := strings.Split(filepath.Base(files[i]), ".")[0]
		hashInMetadata, err := es.HasHash(hash)
		if err != nil {
			log.Println(err)
			return
		}
		if !hashInMetadata {
			del(hash)
		}
	}
}

func del(hash string) {
	log.Println("delete ", hash)
	addr := "http://" + config.Configs.ServerAddress + "/objects/" + hash
	request, _ := http.NewRequest("DELETE", addr, nil)
	client := http.Client{}
	client.Do(request)
}
