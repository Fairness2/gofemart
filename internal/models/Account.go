package models

import (
	"database/sql"
	"time"
)

type Account struct {
	ID          int64          `db:"id"`
	UserID      int64          `db:"user_id"`
	Difference  int            `db:"difference"`
	OrderNumber sql.NullString `db:"order_number"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

func NewAccount(orderNumber sql.NullString, userID int64, difference int) *Account {
	return &Account{
		Difference:  difference,
		UserID:      userID,
		OrderNumber: orderNumber,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

type Balance struct {
	Current   int `db:"current" json:"current"`
	Withdrawn int `db:"withdrawn" json:"withdrawn"`
}
