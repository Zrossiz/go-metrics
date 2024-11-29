// support package for check body hash
package hashgenerator

import (
	"crypto/sha256"
	"encoding/hex"
)

func Generate(body []byte, key string) string {
	if key == "" {
		return ""
	}
	h := sha256.New()
	h.Write(body)
	h.Write([]byte(key))
	generatedHash := hex.EncodeToString(h.Sum(nil))
	return generatedHash
}
