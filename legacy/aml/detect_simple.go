package aml

import (
	"slices"
)

// DetectSuspiciousUsers use an array or transaction to determine suspicious users,
// users who have sent more than 10k euros on a 24h windows are considered suspicious.
// @param slices of transaction to verify
// @return array of suspicious userId
func DetectSuspiciousUsers(transactions []Transaction) []string {
	sortedTx := slices.Clone(transactions)
	slices.SortFunc(sortedTx, func(a, b Transaction) int {
		if a.Timestamp.Before(b.Timestamp) {
			return -1
		} else if a.Timestamp.After(b.Timestamp) {
			return 1
		}
		return 0
	})
	groupedPerUser := groupTxPerUsers(sortedTx)

	suspicious := make([]string, 0)
	for k, v := range groupedPerUser {
		if checkMaxAmountOnRollingWindow(MaxAmount, v) {
			suspicious = append(suspicious, k)
		}
	}

	return suspicious
}

func checkMaxAmountOnRollingWindow(maxAmount float64, txs []Transaction) bool {

	// for each data point we will verify the next transaction in a N hours rolling window.
	for i, tx := range txs {
		total := tx.Amount
		start := tx.Timestamp

		for _, ts := range txs[i+1:] {
			if ts.Timestamp.Sub(start) > SuspiciousWindowSize {
				break
			}
			total += ts.Amount
		}
		if total >= maxAmount {
			return true
		}
	}
	return false
}

func groupTxPerUsers(transactions []Transaction) map[string][]Transaction {
	groupedTransaction := make(map[string][]Transaction)

	for _, t := range transactions {
		groupedTransaction[t.UserID] = append(groupedTransaction[t.UserID], t)
	}

	return groupedTransaction
}
