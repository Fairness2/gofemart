-- +goose Up
create table public.d_order_status
(
    code        varchar(10)
        constraint d_order_status_pk
            primary key,
    description varchar
);
comment on column public.d_order_status.code is 'Код статуса';
comment on column public.d_order_status.description is 'Описание статуса';
INSERT INTO d_order_status (code, description) VALUES ('NEW', 'Заказ загружен в систему, но не попал в обработку');
INSERT INTO d_order_status (code, description) VALUES ('PROCESSING', 'Вознаграждение за заказ рассчитывается');
INSERT INTO d_order_status (code, description) VALUES ('INVALID', 'Система расчёта вознаграждений отказала в расчёте');
INSERT INTO d_order_status (code, description) VALUES ('PROCESSED', 'Данные по заказу проверены и информация о расчёте успешно получена');
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
    last_checked_at timestamp
);
comment on table public.t_order is 'Заказы пользователя';
comment on column public.t_order.number is 'Цифровой код заказа';
comment on column public.t_order.user_id is 'Пользователь';
comment on column public.t_order.created_at is 'Дата создания записи';
comment on column public.t_order.updated_at is 'Дата обновления записи';
create index t_order_user_id_index on public.t_order (user_id);

-- +goose Down
