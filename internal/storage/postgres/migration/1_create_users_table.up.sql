create table if not exists users
(
    id               bigint unique not null primary key,
    username         varchar,
    first_name       varchar,
    last_name        varchar,
    language_code    varchar,
    wallet           varchar,
    first_message_at timestamp,
    last_message_at  timestamp
);