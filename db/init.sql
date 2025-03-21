BEGIN;

CREATE TABLE IF NOT EXISTS memes
(
    id VARCHAR(63) PRIMARY KEY,
    board_id VARCHAR(63),
    filename TEXT,
    descriptions TEXT
);

CREATE TABLE IF NOT EXISTS boards
(
    id VARCHAR(63) PRIMARY KEY,
    owner_id VARCHAR(63),
    name TEXT
);

CREATE TABLE IF NOT EXISTS users
(
    id VARCHAR(63) PRIMARY KEY,
    login TEXT,
    password TEXT
);

CREATE TABLE IF NOT EXISTS medias
(
    id VARCHAR(63) PRIMARY KEY,
    body BYTEA
);

COMMIT;