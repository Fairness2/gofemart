package repositories

import (
	"context"
	"database/sql"
	"errors"
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
func (r *OrderRepository) GetOrdersExcludeOrdersWhereStatusIn(limit int, excludedNumbers []string, olderThen time.Time, statuses ...string) ([]models.Order, error) {
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
	wheres = append(wheres, olderThen, olderThen)
	wheres = append(wheres, limit)
	sqlStr := "SELECT * FROM t_order WHERE " + statusSQL + " AND " + excludedNumbersSQL + " AND ((last_checked_at NOT NULL AND last_checked_at <= ?) OR (last_checked_at IS NULL AND created_at <= ?)) LIMIT ?"
	var orders []models.Order
	err = r.db.SelectContext(r.ctx, &orders, sqlStr, wheres...)

	return orders, err
}

func (r *OrderRepository) GetOrderByNumber(number string) (*models.Order, bool, error) {
	var order models.Order
	err := r.db.QueryRowxContext(r.ctx, "SELECT * FROM t_order WHERE number = $1", number).StructScan(&order)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &order, true, nil
}

func (r *OrderRepository) GetOrdersByUserWithAccrual(userId int64) ([]models.OrderWithAccrual, error) {
	var orders []models.OrderWithAccrual
	err := r.db.SelectContext(r.ctx, &orders, "SELECT t.*, ta.difference accrual FROM t_order t LEFT JOIN t_account ta ON t.number = ta.order_number AND ta.difference > 0 WHERE t.user_id = $1", userId)
	return orders, err
}

func (r *OrderRepository) GetOrdersByUserWithdraw(userId int64) ([]models.OrderWithdraw, error) {
	var orders []models.OrderWithdraw
	err := r.db.SelectContext(r.ctx, &orders, "SELECT t.*, abs(ta.difference) accrual, ta.created_at processed_at FROM  t_order t INNER JOIN public.t_account ta on t.number = ta.order_number AND ta.difference < 0 WHERE t.user_id = $1;", userId)
	return orders, err
}
