package ipchecker

import (
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/config"
)

func TrustedIpCheckerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		senderIp := r.Header.Get("X-Real-IP")
		configTrustedSubnet := config.AppConfig.TrustedSubnet

		if configTrustedSubnet == "" {
			next.ServeHTTP(w, r)
		}

		if configTrustedSubnet != senderIp {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
