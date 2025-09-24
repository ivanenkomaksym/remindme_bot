package middleware

import (
	"net/http"
	"strings"

	"github.com/ivanenkomaksym/remindme_bot/bootstrap"
)

// APIKeyMiddleware protects /api endpoints using X-API-Key header
func APIKeyMiddleware(app *bootstrap.Application, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Bypass non-API routes and health/webhook endpoints
		if path == "/" || path == "/telegram-webhook" || !strings.HasPrefix(path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}

		expectedKey := app.Env.Config.App.APIKey
		if expectedKey == "" {
			http.Error(w, "API key not configured", http.StatusInternalServerError)
			return
		}

		providedKey := r.Header.Get("X-API-Key")
		if providedKey == "" || providedKey != expectedKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
