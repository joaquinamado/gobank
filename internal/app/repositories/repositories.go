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

	return &Repositories{
		Account: &accountRepo{
			store: *store,
		},
		Transfer: &transferRepo{
			store: *store,
		},
	}, nil
}
