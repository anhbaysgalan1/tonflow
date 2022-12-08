create table if not exists pictures
(
    id       varchar unique not null primary key,
    added_at timestamp
);