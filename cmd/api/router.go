package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strings"
	"time"
)

func (api *Api) setupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	//r.Use(api.SendRequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	// custom 405 not allowed handler
	r.MethodNotAllowed(api.methodNotAllowedResponse)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/ping", api.PingHandler)
		r.Options("/ping", api.PingHandler)
		r.Get("/health", api.HealthCheckHandler)
	})

	return r
}

func (api *Api) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	rctx := chi.RouteContext(r.Context())
	routePath := r.URL.Path
	if routePath == "" {
		if r.URL.RawPath != "" {
			routePath = r.URL.RawPath
		}
	}
	tctx := chi.NewRouteContext()
	allowedOptions := []string{"OPTIONS"}
	testOptions := []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodHead, http.MethodDelete, http.MethodConnect, http.MethodTrace}
	for _, option := range testOptions {
		if rctx.Routes.Match(tctx, option, routePath) {
			allowedOptions = append(allowedOptions, option)
		}
	}
	w.Header().Set("Allow", strings.Join(allowedOptions, ","))
	api.respondError(w, http.StatusMethodNotAllowed, message)
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
