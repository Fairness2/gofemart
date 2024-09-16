package repositories

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

// SQLExecutor интерфейс с нужными функциями из sqlx.DB
type SQLExecutor interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}
