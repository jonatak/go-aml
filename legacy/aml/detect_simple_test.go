package aml

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func parse(datetime string) time.Time {
	val, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		log.Fatalf("Invalid datetime in test %v.", err)
	}
	return val
}

func TestDetectSuspiciousUsers(t *testing.T) {

	cases := []struct {
		name   string
		txns   []Transaction
		expect []string
	}{
		{
			name: "Simple suspicious case",
			txns: []Transaction{
				{ID: "1", UserID: "u1", Amount: 4000, Timestamp: parse("2025-07-01T10:00:00Z")},
				{ID: "2", UserID: "u1", Amount: 7000, Timestamp: parse("2025-07-01T12:00:00Z")},  // outside 24h
				{ID: "3", UserID: "u1", Amount: 7000, Timestamp: parse("2025-07-02T01:00:00Z")},  // within 24h → suspicious
				{ID: "4", UserID: "u2", Amount: 10000, Timestamp: parse("2025-07-01T11:00:00Z")}, // edge case
				{ID: "5", UserID: "u3", Amount: 11000, Timestamp: parse("2025-07-01T12:00:00Z")}, // single txn suspicious
				{ID: "6", UserID: "u4", Amount: 7000, Timestamp: parse("2025-07-01T15:00:00Z")},  // cleared
			},
			expect: []string{"u1", "u2", "u3"},
		},
		{
			name: "No suspicious activity",
			txns: []Transaction{
				{ID: "1", UserID: "u1", Amount: 4000, Timestamp: parse("2025-07-01T10:00:00Z")},
				{ID: "4", UserID: "u2", Amount: 500, Timestamp: parse("2025-07-01T11:00:00Z")},
				{ID: "5", UserID: "u3", Amount: 1000, Timestamp: parse("2025-07-01T12:00:00Z")},
			},
			expect: []string{},
		},
		{
			name:   "Empty transaction",
			txns:   []Transaction{},
			expect: []string{},
		},
		{
			name: "Out of order transaction",
			txns: []Transaction{
				{ID: "1", UserID: "u1", Amount: 4000, Timestamp: parse("2025-07-01T10:00:00Z")},
				{ID: "2", UserID: "u1", Amount: 7000, Timestamp: parse("2025-07-02T09:00:00Z")},
				{ID: "3", UserID: "u1", Amount: 7000, Timestamp: parse("2025-07-01T15:00:00Z")},
			},
			expect: []string{"u1"},
		},
		{
			name: "Overlapping windows",
			txns: []Transaction{
				{ID: "1", UserID: "u1", Amount: 5000, Timestamp: parse("2025-07-01T01:00:00Z")},
				{ID: "2", UserID: "u1", Amount: 4000, Timestamp: parse("2025-07-02T00:30:00Z")},
				{ID: "3", UserID: "u1", Amount: 4000, Timestamp: parse("2025-07-03T00:00:00Z")},
				{ID: "4", UserID: "u1", Amount: 4000, Timestamp: parse("2025-07-03T23:59:00Z")},
			},
			expect: []string{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := DetectSuspiciousUsers(c.txns)
			assert.ElementsMatchf(t, got, c.expect, "Expected %v and wanted %v", c.expect, got)
		})
	}

}
