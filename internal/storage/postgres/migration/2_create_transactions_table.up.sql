create table if not exists transactions
(
    source  varchar        not null,
    hash    varchar unique not null,
    value   integer        not null,
    comment varchar        not null
);