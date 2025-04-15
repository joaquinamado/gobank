package types

import (
	"golang.org/x/crypto/bcrypt"
)

func (a *Account) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(password))
}
