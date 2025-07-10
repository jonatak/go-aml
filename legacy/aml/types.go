package aml

import (
	"jonatak/aml/proto"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Transaction struct {
	ID        string
	UserID    string
	Amount    float64
	Timestamp time.Time
}

func (t Transaction) ToPbTransaction() *proto.Transaction {
	return &proto.Transaction{
		Id:        t.ID,
		UserId:    t.UserID,
		Amount:    t.Amount,
		Timestamp: timestamppb.New(t.Timestamp),
	}
}

type TransactionQuery struct {
	Tx           *proto.Transaction
	ResponseChan chan proto.TransactionStatus
}
