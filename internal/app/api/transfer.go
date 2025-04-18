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

// @Summary		Transfer
// @Description	Transfer money between accounts
// @Tags			transfer
// @Accept			json
// @Param			id		path	int	true	"Account ID"
// @Param			page	query   int	false	"Page number"
// @Param			size	query	int	false	"Page size"
// @Security		BearerAuth
// @Success		200
// @Router			/transfer/{id} [get]
func (s *APIServer) handleGetTransferById(w http.ResponseWriter, r *http.Request) error {
	id, err := getPathIntParam(r, "id")

	if err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
	}
	page, errPage := getQueryIntParam(r, "page")
	size, errSize := getQueryIntParam(r, "size")

	query := new(types.PaginationQuery)

	query.Id = id
	if errPage == nil {
		query.Page = page
	}

	if errSize == nil {
		query.Size = size
	}

	transfers, err := s.repo.Transfer.GetTransfers(query)

	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}

	return WriteJson(w, http.StatusOK, transfers)

}
