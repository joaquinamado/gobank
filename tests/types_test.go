package main

import (
	"fmt"
	"testing"

	"github.com/joaquinamado/gobank/internal/app/repositories"
	"github.com/stretchr/testify/assert"
)

var rep, err = repositories.NewRepository()

func TestNewAccount(t *testing.T) {
	if err != nil {
		return
	}
	account, err := rep.Account.CreateAccount("John", "Doe", "Pass.1234")
	assert.Nil(t, err)

	fmt.Printf("%+v\n", account)
}
