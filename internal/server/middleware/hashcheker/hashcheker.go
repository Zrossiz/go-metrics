package hashcheker

import (
	"bytes"
	"io"
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/config"
	"github.com/Zrossiz/go-metrics/internal/server/libs/hashgenerator"
)

func HashCheker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hash := r.Header.Get("HashSHA256")
		if hash != "" {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "failed to read request body", http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			generatedHash := hashgenerator.Generate(bodyBytes, config.AppConfig.ApiKey)

			if generatedHash == hash {
				next.ServeHTTP(w, r)
			}

			http.Error(w, "invalid hash", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
