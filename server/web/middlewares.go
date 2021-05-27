package web

import (
	"net/http"
	"time"

	"oss.navercorp.com/metis/metis-server/internal/log"
)

func elapsedTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		next.ServeHTTP(writer, request)
		log.Logger.Infof("WEB : /auth %s", time.Since(start))
	})
}
