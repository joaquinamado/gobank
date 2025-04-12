package api

import (
	"net/http"
)

// @Summary		Health Check
// @Description	Check if the server is up and running
// @Tags			health
// @Success		200
// @Router			/health [get]
func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) error {
	return WriteJson(w, http.StatusOK, "OK")
}
