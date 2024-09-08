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
			&migrator.Migration{
				Name: "Create user table",
				Func: func(tx *sql.Tx) error {
					if _, err := tx.Exec(`
create table public.d_order_status
(
    code        varchar(10)
        constraint d_order_status_pk
            primary key,
    description varchar
);
`); err != nil {
						return err
					}
					if _, err := tx.Exec("comment on column public.d_order_status.code is 'Код статуса';"); err != nil {
						return err
					}
					if _, err := tx.Exec("comment on column public.d_order_status.description is 'Описание статуса';"); err != nil {
						return err
					}
					if _, err := tx.Exec("INSERT INTO d_order_status (code, description) VALUES ('NEW', 'Заказ загружен в систему, но не попал в обработку');"); err != nil {
						return err
					}
					if _, err := tx.Exec("INSERT INTO d_order_status (code, description) VALUES ('PROCESSING', 'Вознаграждение за заказ рассчитывается');"); err != nil {
						return err
					}
					if _, err := tx.Exec("INSERT INTO d_order_status (code, description) VALUES ('INVALID', 'Система расчёта вознаграждений отказала в расчёте');"); err != nil {
						return err
					}
					if _, err := tx.Exec("INSERT INTO d_order_status (code, description) VALUES ('PROCESSED', 'Данные по заказу проверены и информация о расчёте успешно получена');"); err != nil {
						return err
					}

					if _, err := tx.Exec(`
create table public.t_order
(
    number          varchar
        constraint t_order_pk
            primary key,
    user_id         integer                   not null
        constraint t_order_t_user_id_fk
            references public.t_user,
    created_at      timestamp   default now() not null,
    updated_at      timestamp   default now() not null,
    status_code     varchar(10) default 'NEW' not null
        constraint t_order_d_order_status_code_fk
            references public.d_order_status (code),
    last_checked_at timestamp,
	accrual_sum integer
);
`); err != nil {
						return err
					}
					if _, err := tx.Exec("comment on table public.t_order is 'Заказы пользователя';"); err != nil {
						return err
					}
					if _, err := tx.Exec("comment on column public.t_order.number is 'Цифровой код заказа';"); err != nil {
						return err
					}
					if _, err := tx.Exec("comment on column public.t_order.user_id is 'Пользователь';"); err != nil {
						return err
					}
					if _, err := tx.Exec("comment on column public.t_order.created_at is 'Дата создания записи';"); err != nil {
						return err
					}
					if _, err := tx.Exec("comment on column public.t_order.updated_at is 'Дата обновления записи';"); err != nil {
						return err
					}
					if _, err := tx.Exec("comment on column public.t_order.status_code is 'Статус заказа';"); err != nil {
						return err
					}
					if _, err := tx.Exec("comment on column public.t_order.accrual_sum is 'Кол-во бонусов';"); err != nil {
						return err
					}
					if _, err := tx.Exec("create index t_order_user_id_index on public.t_order (user_id);"); err != nil {
						return err
					}

					return nil
				},
			},
			&migrator.Migration{
				Name: "Create account table",
				Func: func(tx *sql.Tx) error {
					if _, err := tx.Exec(`create table if not exists t_account
(
    id           bigserial
        constraint t_account_pk
            primary key,
    difference   integer not null,
    user_id      bigint  not null
        constraint t_account_t_user_id_fk
            references t_user,
    order_number varchar
        constraint t_account_t_order_number_fk
            references t_order.
	created_at      timestamp   default now() not null,
    updated_at      timestamp   default now() not null,
);`); err != nil {
						return err
					}
					if _, err := tx.Exec("comment on table t_account is 'Счёт бонусов пользователей';"); err != nil {
						return err
					}
					if _, err := tx.Exec("comment on column t_account.id is 'Идентификатор транзакции';"); err != nil {
						return err
					}
					if _, err := tx.Exec("comment on column t_account.difference is 'На счколько изменился счёт';"); err != nil {
						return err
					}
					if _, err := tx.Exec("comment on column t_account.user_id is 'Какому пользователю принадлежит';"); err != nil {
						return err
					}
					if _, err := tx.Exec("comment on column t_account.order_number is 'Идентификатор связанного заказа';"); err != nil {
						return err
					}
					if _, err := tx.Exec("create index if not exists t_account_user_id_index on t_account (user_id);"); err != nil {
						return err
					}
					return nil
				},
			},
		),
	)

	return m, err
}
