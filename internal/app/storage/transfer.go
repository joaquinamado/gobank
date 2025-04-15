package storage

import (
	"database/sql"
	"fmt"

	"github.com/joaquinamado/gobank/internal/app/types"
)

type transfer interface {
	createTransferTable() error
	CreateTransfer(*types.Transfer) error
}

type postgresTransfer struct {
	db *sql.DB
}

func (s *postgresTransfer) createTransferTable() error {
	query := `create table if not exists transfer (
              id serial primary key,
			  sender_id int,
			  receiver_id int,
			  FOREIGN KEY (sender_id) REFERENCES account (id),
			  FOREIGN KEY (receiver_id) REFERENCES account (id),
              amount serial,
              created_at timestamp
            )`

	_, err := s.db.Exec(query)
	return err
}

func (s *postgresTransfer) CreateTransfer(trans *types.Transfer) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	_, err = tx.Exec(`
		INSERT INTO transfer (
			sender_id, receiver_id, amount, created_at
		) VALUES ($1, $2, $3, $4)
	`, trans.SenderId, trans.ReceiverId, trans.Amount, trans.CreatedAt)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE account SET balance = balance - $1 WHERE id = $2
	`, trans.Amount, trans.SenderId)
	if err != nil {
		return err
	}
	fmt.Println("LLEGA 3")

	_, err = tx.Exec(`
		UPDATE account SET balance = balance + $1 WHERE id = $2
	`, trans.Amount, trans.ReceiverId)
	if err != nil {
		return err
	}

	return nil
}
