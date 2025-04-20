package storage

import (
	"database/sql"
	"fmt"

	"github.com/joaquinamado/gobank/internal/app/types"
)

type transfer interface {
	createTransferTable() error
	CreateTransfer(*types.Transfer) error
	GetTransfers(*types.PaginationQuery) ([]*types.Transfer, error)
	GetInvoiceInformation(int) (*types.TransferInvoice, error)
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

	_, err = tx.Exec(`
		UPDATE account SET balance = balance + $1 WHERE id = $2
	`, trans.Amount, trans.ReceiverId)
	if err != nil {
		return err
	}

	return nil
}

func (s *postgresTransfer) GetTransfers(query *types.PaginationQuery) ([]*types.Transfer, error) {
	sqlQuery := `
		SELECT * FROM transfer 
		WHERE sender_id = $1
		ORDER BY created_at DESC
		LIMIT $2
		OFFSET $3
	`
	rows, err := s.db.Query(sqlQuery,
		query.Id,
		query.Size,
		query.Page,
	)

	if err != nil {
		return nil, err
	}

	transfers := []*types.Transfer{}
	for rows.Next() {
		transfer, err := scanIntoTransfer(rows)

		if err != nil {
			return nil, err
		}

		transfers = append(transfers, transfer)
	}

	return transfers, nil

}

func (s *postgresTransfer) GetInvoiceInformation(id int) (*types.TransferInvoice, error) {
	query := ` 
		SELECT t.id, t.amount, t.created_at,     
		RTRIM(LTRIM(
				CONCAT(
					COALESCE(a.first_name || ' ', '')
					, COALESCE(a.last_name, '')
				)
			)) AS sender_fullname,
		RTRIM(LTRIM(
				CONCAT(
					COALESCE(a2.first_name || ' ', '')
					, COALESCE(a2.last_name, '')
				)
			)) AS receiver_fullname,
		a."number" as sender_account,
		a2."number" as receiver_account
		FROM transfer t 
		LEFT JOIN account a ON  t.sender_id = a.id
		LEFT JOIN account a2 ON  t.receiver_id = a2.id
		WHERE t.id = $1;
	`

	rows, err := s.db.Query(query, id)

	if err != nil {
		return nil, err
	}

	if rows.Next() {
		invoice := new(types.TransferInvoice)
		err = rows.Scan(
			&invoice.ID,
			&invoice.Amount,
			&invoice.CreatedAt,
			&invoice.SenderFullName,
			&invoice.ReceiverFullName,
			&invoice.SenderAccount,
			&invoice.ReceiverAccount,
		)

		if err != nil {
			return nil, err
		}

		return invoice, nil
	}

	return nil, fmt.Errorf("Error - invoice not found")

}

func scanIntoTransfer(rows *sql.Rows) (*types.Transfer, error) {
	transfer := new(types.Transfer)
	err := rows.Scan(
		&transfer.ID,
		&transfer.ReceiverId,
		&transfer.SenderId,
		&transfer.Amount,
		&transfer.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return transfer, nil
}
