package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	types "github.com/joaquinamado/gobank/internal/app/types"
)

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("Method not allowed: %s", r.Method)
}

// @Summary		Account
// @Description	Get all accounts
// @Tags			account
// @Success		200
// @Router			/account [get]
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.repo.Account.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, accounts)
}

// @Summary		Account
// @Description	Get account by ID
// @Tags			account
// @Param			id				path	int		true	"Account ID"
// @Success		200
// @Router			/account/{id} [get]
func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	} else {

		id, err := getId(r)

		if err != nil {
			return err
		}

		account, err := s.repo.Account.GetAccountByID(id)

		if err != nil {
			return err
		}

		return WriteJson(w, http.StatusOK, account)
	}
}

// @Summary		Account
// @Description	Create an account
// @Tags			account
// @Param			Data body	types.CreateAccountRequest true	"Create Account Data"
// @Success		200
// @Router			/account [post]
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(types.CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	if err := Validate.Struct(req); err != nil {
		return err
	}

	account, err := s.repo.Account.CreateAccount(req.FirstName, req.LastName, req.Password)

	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, account)
}

// @Summary		Account
// @Description	Update an account
// @Tags			account
// @Param			Data body	types.UpdateAccountRequest true	"Update Account Data"
// @Success		200
// @Router			/account [put]
func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(types.UpdateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	if err := Validate.Struct(req); err != nil {
		return err
	}

	account, err := s.repo.Account.UpdateAccount(req)

	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, account)
}

// @Summary		Account
// @Description	Delete an account
// @Tags			account
// @Param			id	path	int	true	"Account ID"
// @Success		200
// @Router			/account/{id} [delete]
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)

	if err != nil {
		return err
	}
	if err := s.repo.Account.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, map[string]int{"deleted": id})
}
