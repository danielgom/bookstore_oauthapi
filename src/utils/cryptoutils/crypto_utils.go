package cryptoutils

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMd5(in string) string {
	hash := md5.New()
	defer hash.Reset()
	hash.Write([]byte(in))
	return hex.EncodeToString(hash.Sum(nil))
}
