package main

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *Api) setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(api.loggingMiddleware, mux.CORSMethodMiddleware(r))

	r.HandleFunc("/ping", api.PingHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/health", api.HealthCheckHandler).Methods(http.MethodGet)

	return r
}

func (api *Api) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	j := simplejson.New()
	j.Set("alive", "true")
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
