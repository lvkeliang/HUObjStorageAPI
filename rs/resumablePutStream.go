package rs

import (
	"HUObjStorageAPI/objectstream"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type resumabletoken struct {
	Name    string
	Size    int64
	Hash    string
	Servers []string
	Uuids   []string
}

type RSResumablePutStream struct {
	*RSPutStream
	*resumabletoken
}

func NewRSResumablePutStream(dataServers []string, name, hash string, size int64) (*RSResumablePutStream, error) {
	putStream, err := NewRSPutStream(dataServers, hash, size)
	if err != nil {
		return nil, err
	}

	uuids := make([]string, ALL_SHARDS)

	for i := range uuids {
		uuids[i] = putStream.writers[i].(*objectstream.TempPutStream).Uuid
	}
	token := &resumabletoken{
		Name:    name,
		Size:    size,
		Hash:    hash,
		Servers: dataServers,
		Uuids:   uuids,
	}

	return &RSResumablePutStream{putStream, token}, nil
}

func (s *RSResumablePutStream) ToToken() string {
	marshaled, _ := json.Marshal(s)
	return base64.StdEncoding.EncodeToString(marshaled)
}

// CurrentSize 获取第一个临时分片大小并乘以DATA_SHARDS返回
func (s *RSResumablePutStream) CurrentSize() int64 {
	res, err := http.Head(fmt.Sprintf("http://%s/temp/%s", s.Servers[0], s.Uuids[0]))
	if err != nil {
		log.Println(err)
		return -1
	}
	if res.StatusCode != http.StatusOK {
		log.Println(res.StatusCode)
		return -1
	}

	size, err := strconv.ParseInt(res.Header.Get("content-length"), 0, 64)
	if err != nil {
		log.Println(err)
		return -1
	}

	sizeAll := size * DATA_SHARDS
	if sizeAll > s.Size {
		sizeAll = s.Size
	}
	return sizeAll
}

func NewRSResumablePutStreamFromToken(token string) (*RSResumablePutStream, error) {
	b, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	var rt resumabletoken
	err = json.Unmarshal(b, &rt)
	if err != nil {
		return nil, err
	}

	writers := make([]io.Writer, ALL_SHARDS)
	for i := range writers {
		writers[i] = &objectstream.TempPutStream{Server: rt.Servers[i], Uuid: rt.Uuids[i]}
	}

	enc := NewEncoder(writers)
	return &RSResumablePutStream{&RSPutStream{enc}, &rt}, nil

}
