package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"gofemart/internal/models"
	"time"
)

// OrderRepository представляет собой хранилище для работы с заказами в базе данных.
type OrderRepository struct {
	// db пул соединений с базой данных, которыми может пользоваться хранилище
	db SQLExecutor
	// storeCtx контекст, который отвечает за запросы
	ctx context.Context
}

// NewOrderRepository создаёт и возвращает новый экземпляр OrderRepository с предоставленным контекстом и интерфейсом выполнения SQL-запросов.
func NewOrderRepository(ctx context.Context, db SQLExecutor) *OrderRepository { // TODO заменить на интерфейс
	return &OrderRepository{
		ctx: ctx,
		db:  db,
	}
}

// CreateOrder вставляем новый заказ
func (r *OrderRepository) CreateOrder(order *models.Order) error {
	_, err := r.db.NamedExecContext(r.ctx, createOrderSQL, order)
	return err
}

// UpdateOrder обновляем существующий заказ
func (r *OrderRepository) UpdateOrder(order *models.Order) error {
	order.UpdatedAt = time.Now()
	_, err := r.db.NamedExecContext(r.ctx, updateOrderSQL, order)
	return err
}

// GetOrdersExcludeOrdersWhereStatusIn получаем заказы с определёнными статусами
func (r *OrderRepository) GetOrdersExcludeOrdersWhereStatusIn(limit int, excludedNumbers []string, olderThen time.Time, statuses ...string) ([]models.Order, error) {
	//wheres := make([]interface{}, 0, len(excludedNumbers)+len(statuses)+3)
	wheres := make([]interface{}, 0, 5)
	var sqlStr string
	var err error
	if len(excludedNumbers) > 0 {
		wheres = append(wheres, statuses, excludedNumbers, olderThen, olderThen, limit)
		sqlStr, wheres, err = sqlx.In(getOrdersExcludeOrdersWhereStatusInWithNumbersSQL, wheres...)
	} else {
		wheres = append(wheres, statuses, olderThen, olderThen, limit)
		sqlStr, wheres, err = sqlx.In(getOrdersExcludeOrdersWhereStatusInWithoutNumbersSQL, wheres...)
	}
	if err != nil {
		return []models.Order{}, err
	}
	sqlStr = r.db.Rebind(sqlStr)

	var orders []models.Order
	err = r.db.SelectContext(r.ctx, &orders, sqlStr, wheres...)

	return orders, err
}

// GetOrderByNumber извлекает заказ из базы данных, используя предоставленный номер заказа.
// Он возвращает заказ, логическое значение, указывающее, был ли заказ найден, и ошибку, если таковая произошла во время выполнения.
func (r *OrderRepository) GetOrderByNumber(number string) (*models.Order, bool, error) {
	var order models.Order
	err := r.db.QueryRowxContext(r.ctx, getOrderByNumberSQL, number).StructScan(&order)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &order, true, nil
}

// GetOrdersByUserWithAccrual извлекает заказы вместе с их начислением для конкретного пользователя по его идентификатору пользователя.
func (r *OrderRepository) GetOrdersByUserWithAccrual(userID int64) ([]models.OrderWithAccrual, error) {
	var orders []models.OrderWithAccrual
	err := r.db.SelectContext(r.ctx, &orders, getOrdersByUserWithAccrualSQL, userID)
	return orders, err
}

// GetOrdersByUserWithdraw извлекает все записи о снятии средств для данного пользователя по его идентификатору.
func (r *OrderRepository) GetOrdersByUserWithdraw(userID int64) ([]models.OrderWithdraw, error) {
	var orders []models.OrderWithdraw
	err := r.db.SelectContext(r.ctx, &orders, getOrdersByUserWithdrawSQL, userID)
	return orders, err
}
