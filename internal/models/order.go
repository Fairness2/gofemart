package models

import (
	"database/sql"
	"time"
)

type Order struct {
	Number        string       `db:"number"`
	UserId        int64        `db:"user_id"`
	StatusCode    string       `db:"status_code"`
	CreatedAt     time.Time    `db:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at"`
	LastCheckedAt sql.NullTime `db:"last_checked_at"`
}

func NewOrder(number string, userId int64) *Order {
	return &Order{
		Number:        number,
		UserId:        userId,
		StatusCode:    StatusNew,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		LastCheckedAt: sql.NullTime{},
	}
}
