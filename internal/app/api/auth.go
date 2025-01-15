package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/joaquinamado/gobank/internal/app/types"
)

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("Method not allowed: %s", r.Method)
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	acc, err := s.store.GetAccountByNumber(int(req.Number))

	if err != nil {
		return err
	}

	token, err := createJWT(acc)
	if err != nil {
		return err
	}

	resp := LoginResponse{
		Token:  token,
		Number: acc.Number,
	}

	if err != nil {
		return err
	}

	if acc.ValidatePassword(req.Password) != nil {
		return fmt.Errorf("invalid password")
	}

	return WriteJson(w, http.StatusOK, resp)
}
