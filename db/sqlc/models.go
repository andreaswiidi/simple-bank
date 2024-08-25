// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
	"time"
)

type Account struct {
	ID        int64
	Fullname  string
	Username  string
	Balance   int64
	Currency  string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type TransactionHistory struct {
	ID                int64
	AccountID         int64
	Amount            int64
	TransactionType   string
	TransferHistoryID sql.NullInt64
	CreatedAt         time.Time
}

type TransfersHistory struct {
	ID            int64
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
	CreatedAt     time.Time
}
