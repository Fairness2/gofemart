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
	smth, err := r.db.PrepareNamed("INSERT INTO t_user (login, password_hash) VALUES (:login, :password_hash) RETURNING id")
	if err != nil {
		return err
	}
	row := smth.QueryRowxContext(r.ctx, user)
	return row.Scan(&user.Id)
}

func (r *UserRepository) GetUserByLogin(login string) (*models.User, bool, error) {
	var user models.User
	err := r.db.QueryRowxContext(r.ctx, "SELECT id, login, password_hash FROM t_user WHERE login = $1", login).StructScan(&user)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &user, true, nil
}

func (r *UserRepository) GetUserById(id int64) (*models.User, bool, error) {
	var user models.User
	err := r.db.QueryRowxContext(r.ctx, "SELECT id, login, password_hash FROM t_user WHERE id = $1", id).StructScan(&user)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &user, true, nil
}
