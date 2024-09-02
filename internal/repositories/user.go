package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"gofemart/internal/models"
)

// SQLExecutor интерфейс с нужными функциями из sql.DB
type SQLExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type UserRepository struct {
	// db пул соединений с базой данных, которыми может пользоваться хранилище
	db *sqlx.DB
	// storeCtx контекст, который отвечает за запросы
	ctx context.Context
}

var UserR *UserRepository

func NewUserRepository(ctx context.Context, db *sqlx.DB) *UserRepository { // TODO заменить на интерфейс
	return &UserRepository{
		ctx: ctx,
		db:  db,
	}
}

// UserExists проверяем наличие пользователя
func (r *UserRepository) UserExists(login string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(r.ctx, "SELECT true FROM t_user WHERE login = $1", login).Scan(&exists)
	// Если у нас нет записей, то возвращаем false
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CreateUser вставляем нового пользователя и присваиваем ему id
func (r *UserRepository) CreateUser(user *models.User) error {
	res, err := r.db.NamedExecContext(r.ctx, "INSERT INTO t_user (login, password_hash) VALUES (:login, :password_hash)", user)
	if err != nil {
		return err
	}
	user.Id, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}
