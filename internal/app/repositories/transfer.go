package repositories

import (
	"fmt"
	"time"

	"github.com/joaquinamado/gobank/internal/app/storage"
	"github.com/joaquinamado/gobank/internal/app/types"
)

type transfer interface {
	CreateTransfer(*types.TransferRequest, int) (*types.Transfer, error)
}

type transferRepo struct {
	store storage.Storage
	acc   account
}

func (t *transferRepo) CreateTransfer(transferReq *types.TransferRequest, senderId int) (*types.Transfer, error) {

	sender, err := t.acc.GetAccountByID(senderId)
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
		SenderId:   senderId,
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
