package models

import (
	"database/sql"
	"time"
)

// Account представляет собой транзакцию по счету пользователя в системе.
// Он хранит такую информацию, как идентификатор пользователя, сумма транзакции,
// номер связанного заказа и временные метки для создания и обновления.
type Account struct {
	ID          int64          `db:"id"`
	UserID      int64          `db:"user_id"`
	Difference  float64        `db:"difference"`
	OrderNumber sql.NullString `db:"order_number"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

// NewAccount создает новый экземпляр Account с указанным номером заказа, идентификатором пользователя и суммой начисления.
func NewAccount(orderNumber sql.NullString, userID int64, difference float64) *Account {
	return &Account{
		Difference:  difference,
		UserID:      userID,
		OrderNumber: orderNumber,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// Balance представляет собой структуру для хранения текущего и снятого баланса пользователя.
// Current хранит текущий баланс пользователя.
// Withdrawn хранит сумму всех снятых средств пользователя.
type Balance struct {
	Current   float64 `db:"current" json:"current"`
	Withdrawn float64 `db:"withdrawn" json:"withdrawn"`
}
