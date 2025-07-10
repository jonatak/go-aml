# 🕵️‍♂️ AML Suspicious Transactions Detection – Test Exercise

> ⚠️ Note: This repository started as a personal practice project and was not initially intended for sharing. Please excuse the informal structure or incomplete parts.

The excercise I set for myself:

## 🎯 Goal

Implement a function `DetectSuspiciousUsers` that takes a list of transactions and returns a list of user IDs flagged as suspicious.

### Suspicion Criteria:
- Any **single transaction** over a threshold (e.g., `10,000`)
- OR multiple transactions that **cumulatively exceed the threshold** within a **24-hour window**

---

## 🧱 Data Structures

```go
type Transaction struct {
    ID        string
    UserID    string
    Amount    int
    Timestamp time.Time
}
```

## Implementation Overview

This was a preparation project to refresh my Go skills ahead of a technical interview.

### Phase 1: Legacy Proof-of-Concept
- In the `legacy/` folder:
A basic implementation (`detect_simple.go`) with test coverage to detect suspicious activity over a 24-hour window.

### Phase 2: Concurrent System Design

- In `internal/aml/`:
    - `grpcserver.go`: gRPC server receiving transactions and passing them to the detection logic.
    - `user_state.go`: Per-user state tracking with built-in locking.
- In `cmd/`:
    - `server/`: Runs the gRPC server.
    - `client/`: Sends random transactions and prints approval status.
- Both can be built via make build.

### Current Status
- ✅ gRPC server and client implemented
- ✅ State management and locking per user
- ✅ Functional checker logic based on time windows and amount
- ❌ REST API not implemented
- ❌ No message broker or async pipeline
- ❌ Minimal testing (except `detect_simple_test.go` in legacy)
- ❌ Not enough doc/comments