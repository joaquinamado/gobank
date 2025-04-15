package storage

import (
	"database/sql"

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
	query := `BEGIN;
	insert into transfer (
      sender_id, receiver_id, amount, created_at
    )
    values ($1, $2, $3, $4);

	update account SET balance = balance - $3
	where number = $1

	update account SET balance = balance + $3
	where number = $2
	COMMIT;
	`
	_, err := s.db.Exec(query,
		trans.SenderId,
		trans.ReceiverId,
		trans.Amount,
		trans.CreatedAt)

	return err
}
