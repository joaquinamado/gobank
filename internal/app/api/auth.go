package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	types "github.com/joaquinamado/gobank/internal/app/types"
)

// @Summary		Login
// @Description	Login to the API
// @Tags			auth
// @Accept          json
// @Param			data body	types.LoginRequest  true	"Login data"
// @Success		200
// @Router			/login [post]
func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("Method not allowed: %s", r.Method)
	}

	var req types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	if err := Validate.Struct(req); err != nil {
		return err
	}

	fmt.Println("LOGIN REQUEST: ", req)

	acc, err := s.store.GetAccountByNumber(int(req.Number))

	if err != nil {
		return err
	}

	token, err := createJWT(acc)
	if err != nil {
		return err
	}

	resp := types.LoginResponse{
		Token:  token,
		Number: acc.Number,
	}

	if acc.ValidatePassword(req.Password) != nil {
		return fmt.Errorf("invalid password")
	}

	return WriteJson(w, http.StatusOK, resp)
}
