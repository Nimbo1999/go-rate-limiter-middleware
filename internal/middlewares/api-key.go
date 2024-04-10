package middlewares

import (
	"context"
	"net/http"
)

type API_KEY string

const HEADER_API_KEY API_KEY = "API_KEY"

func ApiKeyRetriever(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(string(HEADER_API_KEY))
		ctx := context.WithValue(r.Context(), HEADER_API_KEY, apiKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetApiKeyFromContext(ctx context.Context) string {
	return ctx.Value(HEADER_API_KEY).(string)
}
