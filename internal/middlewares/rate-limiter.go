package middlewares

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nimbo1999/go-rate-limit/internal/services"
	"github.com/nimbo1999/go-rate-limit/internal/utils"
)

type RateLimiterMiddleware = func(http.Handler) http.Handler

func RateLimiter(limiter services.RateLimiterChecker) RateLimiterMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			checkLimiterFunc := CheckLimiter(w, r)
			apiKey := GetApiKeyFromContext(r.Context())
			fmt.Println(apiKey)

			if apiKey != "" {
				log.Println("Request received with API_KEY")
				log.Println("API_KEY:", apiKey)

				status, err := limiter.CheckApiKeyLimit(apiKey)
				if err != nil {
					log.Panicln(err)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("There was an unexpected error on the server"))
					return
				}
				checkLimiterFunc(status, next)
				return
			}

			status, err := limiter.CheckIpLimit(utils.FormattIp(r.RemoteAddr))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("There was an unexpected error on the server"))
				return
			}
			checkLimiterFunc(status, next)
		})
	}
}

func CheckLimiter(w http.ResponseWriter, r *http.Request) func(status int, next http.Handler) {
	return func(status int, next http.Handler) {
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
	}
}

func sendLimitResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
}
