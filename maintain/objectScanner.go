package maintain

import (
	"HUObjStorageAPI/config"
	"HUObjStorageAPI/es"
	"HUObjStorageAPI/objects"
	"HUObjStorageAPI/util"
	"log"
	"path/filepath"
	"strings"
)

// ObjectScanner 检查和修复对象数据
func ObjectScanner() {
	files, _ := filepath.Glob(config.Configs.StorageRoot + "/objects/*")
	for i := range files {
		hash := strings.Split(filepath.Base(files[i]), ".")[0]
		verify(hash)
	}
}

func verify(hash string) {
	log.Println("verify ", hash)
	size, err := es.SearchHashSize(hash)
	if err != nil {
		log.Println(err)
		return
	}

	stream, err := objects.GetStream(hash, size)
	if err != nil {
		log.Println(err)
		return
	}

	d := util.CalculateHash(stream)
	if d != hash {
		log.Printf("object hash mismatch, calculated = %s, requested = %s", d, hash)
	}

	stream.Close()
}
