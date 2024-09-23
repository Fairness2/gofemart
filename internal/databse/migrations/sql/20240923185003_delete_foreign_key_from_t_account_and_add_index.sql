-- +goose Up
create index t_account_order_number_index on public.t_account (order_number);
alter table public.t_account
drop constraint t_account_t_order_number_fk;

-- +goose Down
