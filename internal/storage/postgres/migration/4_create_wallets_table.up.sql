create table if not exists wallets
(
    address    varchar unique not null,
    "version"    integer,
    seed       varchar unique not null,
    created_at timestamp
);