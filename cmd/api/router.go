package main

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func (api *Api) setupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/ping", api.PingHandler)
		r.Options("/ping", api.PingHandler)
		r.Get("/health", api.HealthCheckHandler)
	})

	return r
}

func (api *Api) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	j := simplejson.New()
	j.Set("alive", "true")
	j.Set("environment", api.env)
	j.Set("version", version)
	api.respondJSON(w, http.StatusOK, j)
}

func (api *Api) PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	if r.Method == http.MethodOptions {
		return
	}
	j := simplejson.New()
	j.Set("ping", "pong")
	api.respondJSON(w, http.StatusOK, j)
}

func (api *Api) respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(response)
}

func (api *Api) respondError(w http.ResponseWriter, code int, message string) {
	api.respondJSON(w, code, map[string]string{"error": message})
}
