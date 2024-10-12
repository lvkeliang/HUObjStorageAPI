package rs

import (
	"HUObjStorageAPI/objectstream"
	"io"
)

type RSResumableGetStream struct {
	*decoder
}

func NewRSResumableGetStream(dataServers []string, uuids []string, size int64) (*RSResumableGetStream, error) {
	readers := make([]io.Reader, ALL_SHARDS)

	var err error
	for i := 0; i < ALL_SHARDS; i++ {
		readers[i], err = objectstream.NewTempGetStream(dataServers[i], uuids[i])
		if err != nil {
			return nil, err
		}
	}

	writers := make([]io.Writer, ALL_SHARDS)
	dec := NewDecoder(readers, writers, size)
	return &RSResumableGetStream{dec}, nil
}
