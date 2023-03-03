begin;

create schema if not exists tonflow;

create table if not exists tonflow.users
(
    id            bigint unique not null,
    username      varchar,
    first_name    varchar,
    last_name     varchar,
    language_code varchar,
    wallet        varchar,
    created_at    timestamp     not null
);

create table if not exists tonflow.wallets
(
    address    varchar unique not null,
    version    integer,
    seed       varchar unique not null,
    created_at timestamp      not null
);

commit;