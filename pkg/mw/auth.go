package mw

import (
	"net/http"
	"pocket/pkg/auth"
	"pocket/pkg/response"
)

func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKeyHeader := r.Header.Get("X-API-Key")
		if apiKeyHeader == "" {
			response.UnauthorizedResponse(w, "Missing API Key")
			return
		}
		isKeyValid := auth.IsKeyValid(apiKeyHeader)
		if !isKeyValid {
			response.UnauthorizedResponse(w, "Invalid API Key")
			return
		}
		next.ServeHTTP(w, r)
	})
}
