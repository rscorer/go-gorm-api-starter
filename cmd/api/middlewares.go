package main

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (api *Api) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.log.Info(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

const requestIDHeader = "X-Request-Id"

func (api *Api) SendRequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if w.Header().Get(requestIDHeader) == "" {
			w.Header().Add(
				requestIDHeader,
				middleware.GetReqID(r.Context()),
			)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
