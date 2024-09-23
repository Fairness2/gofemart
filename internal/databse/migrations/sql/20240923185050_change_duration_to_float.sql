-- +goose Up
alter table public.t_account alter column difference type double precision using difference::double precision;

-- +goose Down
