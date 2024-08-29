package migrations

import (
	"database/sql"
	"github.com/lopezator/migrator"
)

func New() (*migrator.Migrator, error) {
	// Configure migrations
	m, err := migrator.New(
		migrator.Migrations(
			&migrator.Migration{
				Name: "Create user table",
				Func: func(tx *sql.Tx) error {
					if _, err := tx.Exec(`create table public.t_user
						(
    						id            bigserial               not null
        						constraint t_user_pk
            					primary key,
    						created_at    timestamp default now() not null,
    						updated_at    timestamp default now() not null,
    						login         varchar                 not null,
    						password_hash varchar                 not null
						);`); err != nil {
						return err
					}
					if _, err := tx.Exec("comment on column public.t_user.id is 'Идентификатор пользователя';"); err != nil {
						return err
					}
					if _, err := tx.Exec("comment on column public.t_user.created_at is 'Дата создания пользователя';"); err != nil {
						return err
					}
					if _, err := tx.Exec("comment on column public.t_user.updated_at is 'Дата обновления пользователя';"); err != nil {
						return err
					}
					if _, err := tx.Exec("create unique index t_user_login_uindex on public.t_user (login);"); err != nil {
						return err
					}
					return nil
				},
			},
		),
	)

	return m, err
}
