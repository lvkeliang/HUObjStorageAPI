package util

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"strings"
)

func CalculateHash(r io.Reader) string {
	hash := sha256.New()
	io.Copy(hash, r)
	// 计算 Base64 编码的哈希值
	base64Hash := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	// 替换 `/` 为 `_`
	return strings.ReplaceAll(base64Hash, "/", "_")
}
