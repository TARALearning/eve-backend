package eve

import (
	"crypto/sha1"
	"encoding/hex"
)

var (
	// SECRETSALT is used to be added as addition to the given word for sha1 hashing
	SECRETSALT = "secret salt /"
)

// Sha1 hashes the given word with the SECRETSALT as content and returns the hashed string
func Sha1(word string) string {
	hash := sha1.New()
	hash.Write([]byte(SECRETSALT + word))
	return hex.EncodeToString(hash.Sum(nil))
}
