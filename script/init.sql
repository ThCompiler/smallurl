CREATE TABLE IF NOT EXISTS urls
(
    id        bigserial   not null primary key,
    long_url  text        not null,
    short_url varchar(10) not null unique
);

CREATE INDEX short_urls on urls (short_url);