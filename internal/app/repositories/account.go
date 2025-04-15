package repositories

import (
	"fmt"
	"log"

	"github.com/joaquinamado/gobank/internal/app/storage"
	types "github.com/joaquinamado/gobank/internal/app/types"

	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type account interface {
	CreateAccount(firstName, lastName, password string) (*types.Account, error)
	SeedAccounts()
	DeleteAccount(int) error
	UpdateAccount(*types.UpdateAccountRequest) (*types.Account, error)
	GetAccounts() ([]*types.Account, error)
	GetAccountByID(int) (*types.Account, error)
	GetAccountByNumber(int) (*types.Account, error)
}

type accountRepo struct {
	store storage.Storage
}

func (r *accountRepo) CreateAccount(firstName, lastName, password string) (*types.Account, error) {
	encp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	account := &types.Account{
		FirstName:         firstName,
		LastName:          lastName,
		EncryptedPassword: string(encp),
		Number:            int64(rand.Intn(1000000)),
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	if err := r.store.Account.CreateAccount(account); err != nil {
		return nil, err
	}

	return account, nil

}

func (r *accountRepo) DeleteAccount(id int) error {
	if err := r.store.Account.DeleteAccount(id); err != nil {
		return err
	}
	return nil
}

func (r *accountRepo) UpdateAccount(account *types.UpdateAccountRequest) (*types.Account, error) {
	return r.store.Account.UpdateAccount(account)
}

func (r *accountRepo) GetAccounts() ([]*types.Account, error) {
	return r.store.Account.GetAccounts()
}

func (r *accountRepo) GetAccountByID(id int) (*types.Account, error) {
	return r.store.Account.GetAccountByID(id)
}

func (r *accountRepo) GetAccountByNumber(number int) (*types.Account, error) {
	return r.store.Account.GetAccountByNumber(number)
}

func (r *accountRepo) seedAccount(fname, lname, pw string) *types.Account {
	acc, err := r.CreateAccount(fname, lname, pw)
	if err != nil {
		log.Fatal(err)
	}

	if err := r.store.Account.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}
	fmt.Println("new account created => ", acc.Number)

	return acc
}

func (r *accountRepo) SeedAccounts() {
	r.seedAccount("John", "Doe", "password")
}
