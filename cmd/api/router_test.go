package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApi_Router(t *testing.T) {
	log := logrus.New()
	log.SetReportCaller(true)
	log.SetLevel(logrus.DebugLevel)
	api := Api{
		db:  nil,
		log: log,
	}

	tests := []struct {
		Name             string
		RequestMethod    string
		RequestUrl       string
		RequestFunction  func(w http.ResponseWriter, r *http.Request)
		ExpectedStatus   int
		ExpectedResponse string
	}{
		{
			Name:             "Test Ping",
			RequestMethod:    "GET",
			RequestUrl:       "/ping",
			RequestFunction:  api.PingHandler,
			ExpectedStatus:   http.StatusOK,
			ExpectedResponse: `{"ping":"pong"}`,
		},
		{
			Name:             "PingHandlerOptions",
			RequestMethod:    "OPTIONS",
			RequestUrl:       "/ping",
			RequestFunction:  api.PingHandler,
			ExpectedStatus:   http.StatusOK,
			ExpectedResponse: ``,
		},
		{
			Name:             "HealthCheck",
			RequestMethod:    "GET",
			RequestUrl:       "/health",
			RequestFunction:  api.HealthCheckHandler,
			ExpectedStatus:   http.StatusOK,
			ExpectedResponse: `{"alive":"true"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			req, err := http.NewRequest(test.RequestMethod, test.RequestUrl, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(test.RequestFunction)
			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != test.ExpectedStatus {
				t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
			}

			if rr.Body.String() != test.ExpectedResponse {
				t.Errorf("handler returned unexpected body; got %v wanted %v", rr.Body.String(), test.ExpectedResponse)
			}
		})
	}
}

func TestApi_RespondError(t *testing.T) {
	log := logrus.New()
	log.SetReportCaller(true)
	log.SetLevel(logrus.DebugLevel)
	api := Api{
		db:  nil,
		log: log,
	}
	rr := httptest.NewRecorder()
	api.respondError(rr, http.StatusForbidden, "forbidden")
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusForbidden)
	}
	expected := `{"error":"forbidden"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body; got %v wanted %v", rr.Body.String(), expected)
	}
}

func TestApi_RespondJsonError(t *testing.T) {
	log := logrus.New()
	log.SetReportCaller(true)
	log.SetLevel(logrus.DebugLevel)
	api := Api{
		db:  nil,
		log: log,
	}
	rr := httptest.NewRecorder()
	api.respondJSON(rr, http.StatusForbidden, func() { return })
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusInternalServerError)
	}
	expected := `json: unsupported type: func()`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body; got %v wanted %v", rr.Body.String(), expected)
	}
}
