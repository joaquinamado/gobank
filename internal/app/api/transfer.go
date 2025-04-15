package api

import (
	"encoding/json"
	"net/http"

	types "github.com/joaquinamado/gobank/internal/app/types"
)

// @Summary		Transfer
// @Description	Transfer money between accounts
// @Tags			transfer
// @Accept			json
// @Param			data	body	types.TransferRequest	true	"Transfer Data"
// @Security		BearerAuth
// @Success		200
// @Router			/transfer [post]
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(types.TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()

	transfer, err := s.repo.Transfer.CreateTransfer(transferReq, int(s.accNumber))

	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}

	return WriteJson(w, http.StatusOK, transfer)
}
