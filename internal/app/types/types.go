package types

import (
	"golang.org/x/crypto/bcrypt"
)

func (a *Account) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(password))
}

type TransferRequest struct {
	ToAccount int     `json:"to_account" validate:"required"`
	Amount    float32 `json:"amount" validate:"required"`
}
