package repositories

import (
	"context"
	"database/sql"
	"errors"
	"gofemart/internal/models"
)

// UserRepository представляет собой хранилище для управления данными пользователя.
type UserRepository struct {
	// db пул соединений с базой данных, которыми может пользоваться хранилище
	db SQLExecutor
	// storeCtx контекст, который отвечает за запросы
	ctx context.Context
}

// NewUserRepository initializes and returns a new UserRepository with the given context and SQLExecutor.
func NewUserRepository(ctx context.Context, db SQLExecutor) *UserRepository { // TODO заменить на интерфейс
	return &UserRepository{
		ctx: ctx,
		db:  db,
	}
}

// UserExists проверяем наличие пользователя
func (r *UserRepository) UserExists(login string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(r.ctx, userExistsSQL, login).Scan(&exists)
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
	smth, err := r.db.PrepareNamed(createUserSQL)
	if err != nil {
		return err
	}
	row := smth.QueryRowxContext(r.ctx, user)
	return row.Scan(&user.ID)
}

// GetUserByLogin извлекает пользователя на основе его логина из базы данных.
// Возвращает пользователя, логическое значение, если найдено, и ошибку.
func (r *UserRepository) GetUserByLogin(login string) (*models.User, bool, error) {
	var user models.User
	err := r.db.QueryRowxContext(r.ctx, getUserByLoginSQL, login).StructScan(&user)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &user, true, nil
}

// GetUserByID извлекает пользователя по его уникальному идентификатору из базы данных.
// Возвращает пользователя, логическое значение, указывающее на существование, и ошибку.
func (r *UserRepository) GetUserByID(id int64) (*models.User, bool, error) {
	var user models.User
	err := r.db.QueryRowxContext(r.ctx, getUserByIDSQL, id).StructScan(&user)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &user, true, nil
}
