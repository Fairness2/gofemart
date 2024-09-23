package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Order представляет собой заказ клиента
type Order struct {
	Number        string       `db:"number" json:"number"`
	UserID        int64        `db:"user_id" json:"-"`
	StatusCode    string       `db:"status_code" json:"status"`
	CreatedAt     time.Time    `db:"created_at" json:"-"`
	UpdatedAt     time.Time    `db:"updated_at" json:"updated_at"`
	LastCheckedAt sql.NullTime `db:"last_checked_at" json:"-"`
}

// NewOrder создает и возвращает новый экземпляр Order с начальным статусом StatusNew.
func NewOrder(number string, userID int64) *Order {
	return &Order{
		Number:        number,
		UserID:        userID,
		StatusCode:    StatusNew,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		LastCheckedAt: sql.NullTime{},
	}
}

// JSONTime представляет собой оболочку для типа time стандартной библиотеки,
// которая обеспечивает пользовательскую сериализацию и сканирование JSON.
type JSONTime struct {
	time.Time
}

// Scan преобразует входное значение в тип JSONTime, гарантируя,
// что значение имеет тип time.Time, или возвращает ошибку, если это не так.
func (jt *JSONTime) Scan(src interface{}) error {
	var ok bool
	jt.Time, ok = src.(time.Time)
	if !ok {
		return fmt.Errorf("cannot scan type %T into JSONTime: %v", src, src)
	}
	return nil
}

// MarshalJSON сериализует объект JSONTime в строку JSON, отформатированную в соответствии со стандартом RFC3339.
func (jt JSONTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(jt.Time.Format(time.RFC3339))
}

// OrderWithAccrual объединяет Order с дополнительной информацией о начислении.
// Включает временную метку последнего обновления и сумму начисления.
type OrderWithAccrual struct {
	Order
	UpdatedAt JSONTime `db:"updated_at" json:"updated_at"`
	Accrual   float64  `db:"accrual" json:"accrual,omitempty"`
}

// OrderWithdraw представляет собой запись о заказе с тратами.
// Number содержит номер заказа.
// Accrual содержит сумму начисления по заказу.
// ProcessedAt содержит дату обработки и время вывода.
type OrderWithdraw struct {
	Number      string   `db:"number" json:"order"`
	Accrual     float64  `db:"accrual" json:"sum,omitempty"`
	ProcessedAt JSONTime `db:"processed_at" json:"processed_at"`
}
