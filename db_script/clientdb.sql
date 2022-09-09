DROP DATABASE clientdb;

CREATE DATABASE clientdb;

\connect clientdb;

DROP DATABASE IF EXISTS users;

CREATE TABLE users (
	user_id SERIAL PRIMARY KEY,
	email TEXT,
	passw bytea
);

DROP DATABASE IF EXISTS userstickets;

CREATE TABLE userstickets (
	entry_id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(user_id),
	ticket TEXT
);

DROP DATABASE IF EXISTS reftokens;

CREATE TABLE refresh_tokens (
    entry_id SERIAL PRIMARY KEY,
    token bytea,
    expire_time BIGINT
);