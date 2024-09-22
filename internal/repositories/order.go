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
func (r *OrderRepository) GetOrdersExcludeOrdersWhereStatusIn(limit int, excludedNumbers []string, olderThen time.Time, statuses ...string) ([]models.Order, error) {
	wheres := make([]interface{}, 0, len(excludedNumbers)+len(statuses)+3)
	sqlStr := "SELECT * FROM t_order WHERE "
	// Добавляем статусы к условию выборки
	if len(statuses) > 0 {

		stSQL, vars, err := sqlx.In("status_code IN (?)", statuses)
		if err != nil {
			return []models.Order{}, err
		}
		sqlStr += stSQL + " AND "
		wheres = append(wheres, vars...)
	}

	// добавляем исключающие номера, которые уже находятся в очереди
	if len(excludedNumbers) > 0 {
		stSQL, vars, err := sqlx.In("number NOT IN (?)", excludedNumbers)
		if err != nil {
			return []models.Order{}, err
		}
		sqlStr += stSQL + " AND "
		wheres = append(wheres, vars...)
	}
	// добавляем условия по дате проверки и лимит записей в одной выборке
	sqlStr += "((last_checked_at NOTNULL AND last_checked_at <= ?) OR (last_checked_at IS NULL AND created_at <= ?)) LIMIT ?"
	wheres = append(wheres, olderThen, olderThen, limit)
	sqlStr = r.db.Rebind(sqlStr)

	var orders []models.Order
	err := r.db.SelectContext(r.ctx, &orders, sqlStr, wheres...)

	return orders, err
}

// GetOrderByNumber извлекает заказ из базы данных, используя предоставленный номер заказа.
// Он возвращает заказ, логическое значение, указывающее, был ли заказ найден, и ошибку, если таковая произошла во время выполнения.
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

// GetOrdersByUserWithAccrual извлекает заказы вместе с их начислением для конкретного пользователя по его идентификатору пользователя.
func (r *OrderRepository) GetOrdersByUserWithAccrual(userID int64) ([]models.OrderWithAccrual, error) {
	var orders []models.OrderWithAccrual
	err := r.db.SelectContext(r.ctx, &orders, "SELECT t.*, CASE WHEN ta.difference NOTNULL THEN difference ELSE 0 END accrual FROM t_order t LEFT JOIN t_account ta ON t.number = ta.order_number AND ta.difference > 0 WHERE t.user_id = $1", userID)
	return orders, err
}

// GetOrdersByUserWithdraw извлекает все записи о снятии средств для данного пользователя по его идентификатору.
func (r *OrderRepository) GetOrdersByUserWithdraw(userID int64) ([]models.OrderWithdraw, error) {
	var orders []models.OrderWithdraw
	err := r.db.SelectContext(r.ctx, &orders, "SELECT ta.order_number number, abs(ta.difference) accrual, ta.created_at processed_at FROM  public.t_account ta WHERE ta.user_id = $1 AND ta.difference < 0 AND ta.order_number NOTNULL", userID)
	return orders, err
}
