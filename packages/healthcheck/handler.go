package healthcheck

import (
	"net/http"

	"dietician.local/packages/response"
)

type IHealthCheckHandler interface {
	Liveness(w http.ResponseWriter) error
	Readiness(w http.ResponseWriter) error
}

type healthCheckHandler struct{}

func NewHealthCheckHandler() IHealthCheckHandler {
	return &healthCheckHandler{}
}

func (h *healthCheckHandler) Liveness(w http.ResponseWriter) error {
	if !Liveness() {
		return response.Error(w, http.StatusInternalServerError, "not healthy")
	}

	return response.Success(w, "healthy", nil)
}

func (h *healthCheckHandler) Readiness(w http.ResponseWriter) error {
	readiness := Readiness()
	if !IsConnectionSuccessful(readiness) {
		return response.Error(w, http.StatusInternalServerError, "not healthy")
	}

	return response.Success(w, "healthy", readiness)
}
