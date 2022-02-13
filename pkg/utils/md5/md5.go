package md5

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(raw string) string {
	h := md5.Sum([]byte(raw))
	return hex.EncodeToString(h[:])
}

