-- +goose Up
create table if not exists t_account
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
    references t_order,
    created_at      timestamp   default now() not null,
    updated_at      timestamp   default now() not null
    );
comment on table t_account is 'Счёт бонусов пользователей';
comment on column t_account.id is 'Идентификатор транзакции';
comment on column t_account.difference is 'На сколько изменился счёт';
comment on column t_account.user_id is 'Какому пользователю принадлежит';
comment on column t_account.order_number is 'Идентификатор связанного заказа';
create index if not exists t_account_user_id_index on t_account (user_id);

-- +goose Down
