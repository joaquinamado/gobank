package storage

import (
	"database/sql"
	"fmt"

	"github.com/joaquinamado/gobank/internal/app/env"
	_ "github.com/lib/pq"
)

type Storage struct {
	Account  account
	Transfer transfer
}

func NewPostgresStore() (*Storage, error) {
	dbName := env.GetString("DB_NAME", "gobank")
	dbUser := env.GetString("DB_USER", "postgres")
	dbPassword := env.GetString("DB_PASSWORD", "postgres")
	dbHost := env.GetString("DB_HOST", "localhost")
	dbPort := env.GetString("DB_PORT", "5432")
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", dbUser, dbName, dbPassword, dbHost, dbPort)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Storage{
		Account: &postgresAccount{
			db: db,
		},
		Transfer: &postgresTransfer{
			db: db,
		},
	}, nil
}

func (s *Storage) Init() error {
	err := s.Account.createAccountTable()
	if err != nil {
		return err
	}
	err = s.Transfer.createTransferTable()
	return err
}
