package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Order struct {
	Number        string       `db:"number" json:"number"`
	UserId        int64        `db:"user_id" json:"-"`
	StatusCode    string       `db:"status_code" json:"status"`
	CreatedAt     time.Time    `db:"created_at" json:"-"`
	UpdatedAt     time.Time    `db:"updated_at" json:"updated_at"`
	LastCheckedAt sql.NullTime `db:"last_checked_at" json:"-"`
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

type JSONTime struct {
	time.Time
}

func (jt *JSONTime) Scan(src interface{}) error {
	var ok bool
	jt.Time, ok = src.(time.Time)
	if !ok {
		return fmt.Errorf("cannot scan type %T into JSONTime: %v", src, src)
	}
	return nil
}

func (jt JSONTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(jt.Time.Format(time.RFC3339))
}

type OrderWithAccrual struct {
	Order
	UpdatedAt JSONTime `db:"updated_at" json:"updated_at"`
	Accrual   int      `db:"accrual" json:"accrual,omitempty"`
}

/*func (o OrderWithAccrual) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		OrderWithAccrual
		UpdatedAt string `json:"updated_at"`
	}{
		OrderWithAccrual: o,
		UpdatedAt:        o.UpdatedAt.Format(time.RFC3339),
	})
}*/

type OrderWithdraw struct {
	Number      string   `db:"number" json:"order"`
	Accrual     int      `db:"accrual" json:"sum,omitempty"`
	ProcessedAt JSONTime `db:"processed_at" json:"processed_at"`
}

/*func (o OrderWithdraw) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		OrderWithdraw
		UpdatedAt   string `json:"-"`
		ProcessedAt string `json:"processed_at"`
	}{
		OrderWithdraw: o,
		UpdatedAt:     o.ProcessedAt.Format(time.RFC3339),
	})
}*/
