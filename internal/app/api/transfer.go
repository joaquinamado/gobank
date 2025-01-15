package api

import (
	"encoding/json"
	"net/http"

	. "github.com/joaquinamado/gobank/internal/app/types"
)

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()

	return WriteJson(w, http.StatusOK, transferReq)
}
