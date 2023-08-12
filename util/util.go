package util

import (
	"bytes"
	"encoding/binary"
	"os"
)

// 将int64转bytes，通过大端序写入
func Int64ToBytes(num int64) []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, num)
	return buf.Bytes()
}

func FileIsExit(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
