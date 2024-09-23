package hashgenerator

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/Zrossiz/go-metrics/internal/server/config"
)

func Generate(body []byte, key string) string {
	h := sha256.New()
	h.Write(body)
	h.Write([]byte(config.AppConfig.Key))
	generatedHash := hex.EncodeToString(h.Sum(nil))
	return generatedHash
}
