CREATE TABLE IF NOT EXISTS news
(
    id      SERIAL PRIMARY KEY,
    title   TEXT    NOT NULL,
    content TEXT    NOT NULL,
    created BIGINT  NOT NULL,
    link    varchar NOT NULL
);
