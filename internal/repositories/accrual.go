package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"gofemart/internal/models"
)

type AccountRepository struct {
	// db пул соединений с базой данных, которыми может пользоваться хранилище
	db *sqlx.DB
	// storeCtx контекст, который отвечает за запросы
	ctx context.Context
}

func NewAccountRepository(ctx context.Context, db *sqlx.DB) *AccountRepository { // TODO заменить на интерфейс
	return &AccountRepository{
		ctx: ctx,
		db:  db,
	}
}

// CreateAccount вставляем новую транзакцию на счёт
func (r *AccountRepository) CreateAccount(account *models.Account) error {
	smth, err := r.db.PrepareNamed("INSERT INTO t_account (user_id, difference, order_number, created_at, updated_at) VALUES (:user_id, :difference, :order_number, :created_at, :updated_at) RETURNING id")
	if err != nil {
		return err
	}
	row := smth.QueryRowxContext(r.ctx, account)
	return row.Scan(&account.Id)
}

// GetSum Получаем текущий баланс пользователя
func (r *AccountRepository) GetSum(userId int64) (int, error) {
	var sum int
	row := r.db.QueryRowContext(r.ctx, "SELECT SUM(difference) FROM t_account WHERE user_id = $1", userId)
	if row.Err() != nil {
		return 0, row.Err()
	}
	err := row.Scan(&sum)
	if err != nil {
		return 0, err
	}
	return sum, nil
}

func (r *AccountRepository) GetBalance(userId int64) (*models.Balance, error) { // TODO транзакция для того, чтобы зафиксировать состояние таблицы
	balance := &models.Balance{}
	row := r.db.QueryRowxContext(r.ctx, "SELECT sum(CASE WHEN difference > 0 THEN difference ELSE 0 END) current, sum(CASE WHEN difference < 0 THEN abs(difference) ELSE 0 END) withdrawn FROM t_account WHERE user_id = $1", userId)
	if row.Err() != nil {
		return nil, row.Err()
	}
	if err := row.StructScan(balance); err != nil {
		return nil, err
	}
	return balance, nil
}

func (r *AccountRepository) GetWithdrawByOrder(orderNumber string) (*models.Account, bool, error) {
	account := &models.Account{}
	err := r.db.QueryRowxContext(r.ctx, "SELECT * FROM t_account WHERE order_number = $1 AND difference < 0", orderNumber).
		StructScan(account)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return account, true, nil
}
