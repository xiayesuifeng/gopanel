package auth

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func EncryptionPassword(password string) string {
	md5Data := md5.Sum([]byte(password))
	sha1Data := sha1.Sum([]byte(md5Data[:]))
	return hex.EncodeToString(sha1Data[:])
}
