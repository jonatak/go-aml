package aml

import (
	"fmt"
	"jonatak/aml/proto"
	"slices"
	"sync"
)

type UserState struct {
	sync.Mutex
	transactions []Transaction
}

func NewUserState() *UserState {
	return &UserState{
		transactions: make([]Transaction, 0),
	}
}

// AuthoriseTransaction will verify AML rule to see if that transaction is good to go.
func (u *UserState) AuthoriseTransaction(newTx *proto.Transaction) proto.TransactionStatus {
	u.Lock()
	defer u.Unlock()

	if n := len(u.transactions); n > 0 && newTx.Timestamp.AsTime().Before(u.transactions[n-1].Timestamp) {
		return proto.TransactionStatus_INVALID_DATE
	}

	cleanupIndex := 0
	amount := newTx.Amount

	// check transaction in the SuspiciousWindowSize and clean older transaction.
	fmt.Printf("UserId: %s, current tx: %v\n", newTx.UserId, u.transactions)
	for _, v := range u.transactions {
		if newTx.Timestamp.AsTime().Sub(v.Timestamp) > SuspiciousWindowSize {
			cleanupIndex++
			continue
		}
		amount += v.Amount
		if amount > MaxAmount {
			return proto.TransactionStatus_MAX_AMOUNT_REACH
		}
	}

	if cleanupIndex != 0 {
		u.transactions = slices.Delete(u.transactions, 0, cleanupIndex)
	}

	u.appendNewTx(newTx)

	return proto.TransactionStatus_APPROVED
}

func (u *UserState) appendNewTx(newTx *proto.Transaction) {
	tx := Transaction{
		ID:        newTx.Id,
		UserID:    newTx.UserId,
		Amount:    newTx.Amount,
		Timestamp: newTx.Timestamp.AsTime(),
	}
	u.transactions = append(u.transactions, tx)
}
