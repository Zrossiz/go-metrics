package cryptochecker

import (
	"io"
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/security"
)

type CryptoChecker func([]byte) ([]byte, error)

func DecryptMiddleware(next http.Handler, checkCryptoBody CryptoChecker) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		encryptedBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, "failed to read request body", http.StatusBadRequest)
			return
		}
		decryptedMessage, err := checkCryptoBody(encryptedBody)
		if err != nil {
			http.Error(rw, "failed to decrypt body", http.StatusBadRequest)
			return
		}

		r = r.WithContext(security.NewDecryptedContext(r.Context(), decryptedMessage))

		next.ServeHTTP(rw, r)
	})
}
