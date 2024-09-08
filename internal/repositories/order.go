package repositories

import (
	"context"
	"github.com/jmoiron/sqlx"
	"gofemart/internal/models"
	"time"
)

type OrderRepository struct {
	// db пул соединений с базой данных, которыми может пользоваться хранилище
	db *sqlx.DB
	// storeCtx контекст, который отвечает за запросы
	ctx context.Context
}

var OrderR *OrderRepository

func NewOrderRepository(ctx context.Context, db *sqlx.DB) *OrderRepository { // TODO заменить на интерфейс
	return &OrderRepository{
		ctx: ctx,
		db:  db,
	}
}

// CreateOrder вставляем новый заказ
func (r *OrderRepository) CreateOrder(order *models.Order) error {
	_, err := r.db.NamedExecContext(r.ctx, "INSERT INTO t_order (number, user_id, status_code) VALUES (:number, :user_id, :status_code)", order)
	return err
}

// UpdateOrder обновляем существующий заказ
func (r *OrderRepository) UpdateOrder(order *models.Order) error {
	order.UpdatedAt = time.Now()
	_, err := r.db.NamedExecContext(r.ctx, "UPDATE t_order SET user_id = :user_id, status_code = :status_code, last_checked_at = :last_checked_at, updated_at = :updated_at WHERE number = :number", order)
	return err
}

// GetOrdersExcludeOrdersWhereStatusIn получаем заказы с определёнными статусами
// TODO выглядит не красиво(
func (r *OrderRepository) GetOrdersExcludeOrdersWhereStatusIn(limit int, excludedNumbers []string, statuses ...string) ([]models.Order, error) {
	statusInt := make([]interface{}, 0, len(statuses))
	excludedNumbersInt := make([]interface{}, 0, len(excludedNumbers))
	wheres := make([]interface{}, 0, len(excludedNumbers)+len(statuses)+1)
	for i, v := range statuses {
		statusInt[i] = v
	}
	for i, v := range excludedNumbers {
		excludedNumbersInt[i] = v
	}
	statusSQL, statusVars, err := sqlx.In("status_code IN (?)", statusInt...)
	if err != nil {
		return []models.Order{}, err
	}
	wheres = append(wheres, statusVars...)
	excludedNumbersSQL, numbersVars, err := sqlx.In("number NOT IN (?)", excludedNumbersInt...)
	if err != nil {
		return []models.Order{}, err
	}
	wheres = append(wheres, numbersVars...)
	wheres = append(wheres, limit)
	sql := "SELECT * FROM t_order WHERE " + statusSQL + " AND " + excludedNumbersSQL + " LIMIT ?"
	var orders []models.Order
	err = r.db.SelectContext(r.ctx, &orders, sql, wheres...)

	return orders, err
}
