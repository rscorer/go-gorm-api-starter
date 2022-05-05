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

type Route struct {
	method string
	route  string
}

var Routes = []*Route{}
var testOptions = []string{http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodHead, http.MethodDelete, http.MethodConnect, http.MethodTrace}

func (api *Api) setupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	//r.Use(api.SendRequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.GetHead)

	r.Use(middleware.Timeout(60 * time.Second))

	r.MethodNotAllowed(api.methodNotAllowedResponse)
	r.NotFound(api.notFoundResponse)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/ping", api.PingHandler)
		r.Options("/ping", api.PingHandler)
		r.Get("/health", api.HealthCheckHandler)
	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		Routes = append(Routes, &Route{route: route, method: method})
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}

	return r
}

func (api *Api) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	api.respondError(w, http.StatusNotFound, "resource not found")
}

func (api *Api) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	routePath := r.URL.Path
	if routePath == "" {
		if r.URL.RawPath != "" {
			routePath = r.URL.RawPath
		}
	}
	var allowedOptions []string
	for _, route := range Routes {
		for _, method := range testOptions {
			if route.route == routePath && route.method == method {
				allowedOptions = append(allowedOptions, route.method)
			}
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
