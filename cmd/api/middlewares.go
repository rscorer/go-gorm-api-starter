package main

import (
	"net/http"
)

func (api *Api) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.log.Info(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
