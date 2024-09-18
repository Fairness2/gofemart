package database

import (
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"gofemart/internal/databse/migrations"
	"gofemart/internal/logger"
)

var ErrorEmptyDSN = errors.New("empty dsn")

// DBPool глобальный пул подключений к базе данных для приложения c функцией закрытия
type DBPool struct {
	DBx *sqlx.DB
}

func newPgDBx(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	// Если дсн не передан, то просто возвращаем созданный пул, он не работоспособен
	if dsn == "" {
		return db, nil
	}
	// Сразу проверим работоспособность соединения
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// NewDB инициализация подключения к бд
func NewDB(dsn string) (*DBPool, error) {
	if dsn == "" {
		return nil, ErrorEmptyDSN
	}
	// Создание пула подключений к базе данных для приложения
	db, err := newPgDBx(dsn)
	if err != nil {
		return nil, err
	}
	pool := &DBPool{
		DBx: db,
	}

	return pool, nil
}

func (p *DBPool) Migrate() error {
	logger.Log.Info("Migrate migrations")
	// Применим миграции
	migrator, err := migrations.New()
	if err != nil {
		return err
	}
	return migrator.Migrate(p.DBx.DB)
}

// Close закрытие базы данных
func (p *DBPool) Close() {
	logger.Log.Info("Closing database connection for defer")
	if p.DBx != nil {
		err := p.DBx.Close()
		if err != nil {
			logger.Log.Error(err)
		}
	}
}
