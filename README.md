# AML Suspicious Transactions Detection – Test Exercise

## Goal

Write a function `DetectSuspiciousUsers` that takes a list of transactions and returns a list of user IDs who are considered suspicious. Suspicion criteria:

- One or more transactions over a threshold (e.g. 10,000).
- OR multiple transactions totaling over a threshold (e.g. 10,000) within a 24-hour window.

## Data Structures

### Transaction Struct

```go
type Transaction struct {
    ID        string
    UserID    string
    Amount    int
    Timestamp time.Time
}