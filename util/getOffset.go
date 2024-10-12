package util

import (
	"strconv"
	"strings"
)

// range格式: 以"bytes=<first>-"开头
func GetOffset(byteRange string) int64 {
	if len(byteRange) < 7 {
		return 0
	}
	if byteRange[:6] != "bytes=" {
		return 0
	}
	bytePos := strings.Split(byteRange[6:], "-")
	offset, _ := strconv.ParseInt(bytePos[0], 0, 64)
	return offset
}
