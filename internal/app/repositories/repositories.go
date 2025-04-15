package repositories

import (
	"log"

	"github.com/joaquinamado/gobank/internal/app/storage"
)

type Repositories struct {
	Account  account
	Transfer transfer
}

func NewRepository() (*Repositories, error) {

	store, err := storage.NewPostgresStore()

	if err != nil {
		return nil, err
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	accountRepo := &accountRepo{
		store: *store,
	}

	transferRepo := &transferRepo{
		store: *store,
		acc:   accountRepo,
	}

	return &Repositories{
		Account:  accountRepo,
		Transfer: transferRepo,
	}, nil
}
