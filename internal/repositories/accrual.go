package repositories

import (
	"context"
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
