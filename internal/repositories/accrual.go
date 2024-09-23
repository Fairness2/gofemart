package repositories

import (
	"context"
	"database/sql"
	"errors"
	"gofemart/internal/models"
)

// AccountRepository предоставляет доступ к данным аккаунтов в базе данных.
// Использует контекст для управления запросами и пул соединений с БД.
type AccountRepository struct {
	// db пул соединений с базой данных, которыми может пользоваться хранилище
	db SQLExecutor
	// storeCtx контекст, который отвечает за запросы
	ctx context.Context
}

// NewAccountRepository creates a new instance of AccountRepository with the provided context and SQLExecutor.
func NewAccountRepository(ctx context.Context, db SQLExecutor) *AccountRepository { // TODO заменить на интерфейс
	return &AccountRepository{
		ctx: ctx,
		db:  db,
	}
}

// CreateAccount вставляем новую транзакцию на счёт
func (r *AccountRepository) CreateAccount(account *models.Account) error {
	smth, err := r.db.PrepareNamed(createAccountSQL)
	if err != nil {
		return err
	}
	row := smth.QueryRowxContext(r.ctx, account)
	return row.Scan(&account.ID)
}

// GetSum Получаем текущий баланс пользователя
func (r *AccountRepository) GetSum(userID int64) (float64, error) {
	var sum float64
	row := r.db.QueryRowContext(r.ctx, getSumSQL, userID)
	if row.Err() != nil {
		return 0, row.Err()
	}
	err := row.Scan(&sum)
	if err != nil {
		return 0, err
	}
	return sum, nil
}

// GetBalance рассчитывает и возвращает текущий и снятый баланс для данного пользователя.
func (r *AccountRepository) GetBalance(userID int64) (*models.Balance, error) { // TODO транзакция для того, чтобы зафиксировать состояние таблицы
	balance := &models.Balance{}
	row := r.db.QueryRowxContext(r.ctx, getBalanceSQL, userID)
	if row.Err() != nil {
		return nil, row.Err()
	}
	if err := row.StructScan(balance); err != nil {
		return nil, err
	}
	return balance, nil
}

// GetWithdrawByOrder извлекает запись о снятии средств по номеру заказа.
func (r *AccountRepository) GetWithdrawByOrder(orderNumber string) (*models.Account, bool, error) {
	account := &models.Account{}
	err := r.db.QueryRowxContext(r.ctx, getWithdrawByOrderSQL, orderNumber).
		StructScan(account)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return account, true, nil
}
