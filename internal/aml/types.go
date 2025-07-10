package aml

import (
	"jonatak/aml/proto"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ApiTransactionStatus string

const (
	APPROVED         ApiTransactionStatus = "approved"
	MAX_AMOUNT_REACH ApiTransactionStatus = "max_amount_reach"
	INVALID_DATE     ApiTransactionStatus = "invalid_date"
)

type Transaction struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

type TransactionQuery struct {
	Tx           *proto.Transaction
	ResponseChan chan proto.TransactionStatus
}

type ApiTransactionResponse struct {
	Status ApiTransactionStatus
}

func (t Transaction) ToPbTransaction() *proto.Transaction {
	return &proto.Transaction{
		Id:        t.ID,
		UserId:    t.UserID,
		Amount:    t.Amount,
		Timestamp: timestamppb.New(t.Timestamp),
	}
}
