package middlewares

import (
	"log"
	"net/http"

	"github.com/nimbo1999/go-rate-limit/internal/services"
)

type RateLimiterMiddleware = func(http.Handler) http.Handler

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := services.RateLimiter{}
		status, err := limiter.CheckLimit(r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("There was an unexpected error on the server"))
		}

		switch status {
		case services.IP_REQUEST_LIMIT:
			log.Println("IP_REQUEST_LIMIT...")
			sendLimitResponse(w)
		case services.TOKEN_REQUEST_LIMIT:
			log.Println("TOKEN_REQUEST_LIMIT...")
			sendLimitResponse(w)
		case services.REQUEST_ALLOWED:
			next.ServeHTTP(w, r)
		}
	})
}

func sendLimitResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
}
