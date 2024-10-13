package maintain

import (
	"HUObjStorageAPI/es"
	"log"
)

var MIN_VERSION_COUNT = 5

// DeleteOldMetadata 删除过期元数据
func DeleteOldMetadata() {
	buckets, err := es.SearchVersionStatus(MIN_VERSION_COUNT + 1)
	if err != nil {
		log.Println(err)
		return
	}
	for i := range buckets {
		bucket := buckets[i]
		for v := 0; v < bucket.Doc_count-MIN_VERSION_COUNT; v++ {
			es.DelMetadata(bucket.Key, v+int(bucket.Min_version.Value))
		}
	}
}
