package repositories

import (
	"fmt"
	"time"

	"github.com/joaquinamado/gobank/internal/app/storage"
	"github.com/joaquinamado/gobank/internal/app/types"
)

type transfer interface {
	CreateTransfer(*types.TransferRequest, int) (*types.Transfer, error)
	GetTransfers(*types.PaginationQuery) ([]*types.Transfer, error)
}

type transferRepo struct {
	store storage.Storage
	acc   account
}

func (t *transferRepo) CreateTransfer(transferReq *types.TransferRequest, senderNumber int) (*types.Transfer, error) {

	if senderNumber == transferReq.ToAccount {
		return nil, fmt.Errorf("Cannot send transers to same account number")
	}

	sender, err := t.acc.GetAccountByNumber(senderNumber)
	if err != nil {
		return nil, err
	}

	if sender.Balance-int64(transferReq.Amount) < 0 {
		return nil, fmt.Errorf("Invalid operation check balance")
	}

	receiver, err := t.acc.GetAccountByNumber(transferReq.ToAccount)
	if err != nil {
		return nil, err
	}

	transfer := &types.Transfer{
		SenderId:   sender.ID,
		ReceiverId: receiver.ID,
		Amount:     int64(transferReq.Amount),
		CreatedAt:  time.Now().UTC(),
	}

	err = t.store.Transfer.CreateTransfer(transfer)

	if err != nil {
		return nil, err
	}

	return transfer, nil
}

func (t *transferRepo) GetTransfers(query *types.PaginationQuery) ([]*types.Transfer, error) {
	// Default vaule for ints
	if query.Size == 0 {
		query.Size = 20
	}

	return t.store.Transfer.GetTransfers(query)
}
