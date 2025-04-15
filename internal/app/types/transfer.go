package types

import "time"

type TransferRequest struct {
	ToAccount int     `json:"to_account" validate:"required"`
	Amount    float32 `json:"amount" validate:"required"`
}

type Transfer struct {
	ID         int       `json:"id"`
	SenderId   int       `json:"sender_id"`
	ReceiverId int       `json:"receiver_id"`
	Amount     int64     `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}
