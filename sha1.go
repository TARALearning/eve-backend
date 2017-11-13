package eve

import (
	"crypto/sha1"
	"encoding/hex"
)

var (
	SECRETSALT = "secret salt /"
)

func Sha1(word string) string {
	hash := sha1.New()
	hash.Write([]byte(SECRETSALT + word))
	return hex.EncodeToString(hash.Sum(nil))
}
