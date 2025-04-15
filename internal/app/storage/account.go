package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/joaquinamado/gobank/internal/app/types"
)

type account interface {
	createAccountTable() error
	CreateAccount(*types.Account) error
	DeleteAccount(int) error
	UpdateAccount(*types.UpdateAccountRequest) (*types.Account, error)
	GetAccounts() ([]*types.Account, error)
	GetAccountByID(int) (*types.Account, error)
	GetAccountByNumber(int) (*types.Account, error)
}

type postgresAccount struct {
	db *sql.DB
}

func (s *postgresAccount) createAccountTable() error {
	query := `create table if not exists account (
              id serial primary key,
              first_name varchar(50),
              last_name varchar(50),
              number serial,
              encrypted_password varchar(255),
              balance serial,
              created_at timestamp,
              updated_at timestamp
            )`

	_, err := s.db.Exec(query)
	return err
}

func (s *postgresAccount) CreateAccount(acc *types.Account) error {
	query := ` insert into account (
      first_name, last_name, encrypted_password, number, balance, created_at, updated_at
    )
    values ($1, $2, $3, $4, $5, $6, $7)`

	resp, err := s.db.Query(query,
		acc.FirstName,
		acc.LastName,
		acc.EncryptedPassword,
		acc.Number,
		acc.Balance,
		acc.CreatedAt,
		acc.UpdatedAt)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)
	return nil
}

func (s *postgresAccount) DeleteAccount(id int) error {
	_, err := s.db.Query("delete from account where id = $1", id)
	return err
}

func (s *postgresAccount) UpdateAccount(acc *types.UpdateAccountRequest) (*types.Account, error) {
	query := `
	 update account set first_name = $1, last_name = $2,
	 balance = $3, updated_at = $4
	 where id = $5`

	resp, err := s.db.Query(query,
		acc.FirstName,
		acc.LastName,
		acc.Balance,
		time.Now().UTC(),
		acc.ID,
	)

	if err != nil {
		return nil, err
	}

	account, err := scanIntoAccount(resp)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *postgresAccount) GetAccounts() ([]*types.Account, error) {
	rows, err := s.db.Query("select * from account")

	if err != nil {
		return nil, err
	}

	accounts := []*types.Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *postgresAccount) GetAccountByID(id int) (*types.Account, error) {
	res, err := s.db.Query("select * from account where id = $1", id)

	if err != nil {
		return nil, err
	}
	for res.Next() {
		return scanIntoAccount(res)

	}
	return nil, fmt.Errorf("account %d not found", id)
}

func (s *postgresAccount) GetAccountByNumber(number int) (*types.Account, error) {
	rows, err := s.db.Query("select * from account where number = $1", number)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account with number %d not found", number)
}

func scanIntoAccount(rows *sql.Rows) (*types.Account, error) {
	account := new(types.Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.EncryptedPassword,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return account, nil
}
