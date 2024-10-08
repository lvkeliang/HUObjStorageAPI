package util

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
)

func CalculateHash(r io.Reader) string {
	hash := sha256.New()
	io.Copy(hash, r)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
