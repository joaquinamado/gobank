package types

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Number int64  `json:"number"`
	Token  string `json:"token"`
}

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	EncryptedPassword string    `json:"-"`
	Number            int64     `json:"number"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"created_at"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=100"`
	LastName  string `json:"last_name" validate:"required,min=3,max=100"`
	Password  string `json:"password" validate:"required,min=8"`
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	encp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		EncryptedPassword: string(encp),
		Number:            int64(rand.Intn(1000000)),
		CreatedAt:         time.Now().UTC(),
	}, nil
}

func (a *Account) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(password))
}

type TransferRequest struct {
	ToAccount int     `json:"to_account" validate:"required"`
	Amount    float32 `json:"amount" validate:"required"`
}
