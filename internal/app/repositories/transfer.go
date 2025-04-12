package repositories

import (
	"github.com/joaquinamado/gobank/internal/app/storage"
)

type transfer interface {
}

type transferRepo struct {
	store storage.Storage
}
