CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    created BIGINT NOT NULL,
    link varchar NOT NULL UNIQUE
);

