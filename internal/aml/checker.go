package aml

import (
	"context"
	"jonatak/aml/proto"
	"sync"
)

// AMLChecker handle the transaction validation flow.
type AMLChecker struct {
	users sync.Map
}

func (a *AMLChecker) StartLoop(ctx context.Context, tx chan TransactionQuery) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case txQuery, ok := <-tx:
			if !ok {
				return nil
			}
			go a.handleRequest(txQuery.Tx, txQuery.ResponseChan)
		}
	}
}

func (a *AMLChecker) handleRequest(newTx *proto.Transaction, rspChan chan proto.TransactionStatus) {
	var userState *UserState

	val, _ := a.users.LoadOrStore(newTx.GetUserId(), NewUserState())

	userState = val.(*UserState)

	status := userState.AuthoriseTransaction(newTx)

	rspChan <- status
}
