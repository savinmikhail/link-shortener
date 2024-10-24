create table IF NOT EXISTS short_link (
    id serial primary key,
    short_code char(8),
    original_url varchar(255)
);

create unique index if not exists short_code on short_link (short_code);


