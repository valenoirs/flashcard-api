package middleware

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/valenoirs/flashcard-api/internal/config"
)

func CORS(cfg *config.Config, logger *slog.Logger) func(http.Handler) http.Handler {
	logger.Info("initializing cors middleware...")
	allowedOriginMap := make(map[string]bool)
	for _, origin := range cfg.CORS.AllowedOrigin {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			allowedOriginMap[trimmed] = true
		}
	}

	methodsStr := strings.Join(cfg.CORS.AllowedMethod, ", ")
	headersStr := strings.Join(cfg.CORS.AllowedHeader, ", ")
	credentialsStr := strconv.FormatBool(cfg.CORS.AllowCredential)
	maxAgeStr := strconv.Itoa(cfg.CORS.MaxAge)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			_, isAllowed := allowedOriginMap[origin]

			if origin == "" {
				isAllowed = true
			}

			if isAllowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", methodsStr)
				w.Header().Set("Access-Control-Allow-Headers", headersStr)
				w.Header().Set("Access-Control-Allow-Credentials", credentialsStr)
				w.Header().Set("Access-Control-Max-Age", maxAgeStr)
			}

			if r.Method == http.MethodOptions {
				if isAllowed {
					w.WriteHeader(http.StatusNoContent)
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
