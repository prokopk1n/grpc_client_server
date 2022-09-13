DROP DATABASE IF EXISTS clientdb;
CREATE DATABASE clientdb;

\connect clientdb;

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
	user_id SERIAL PRIMARY KEY,
	email TEXT,
	passw bytea
);

DROP TABLE IF EXISTS userstickets CASCADE;
CREATE TABLE userstickets (
	entry_id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(user_id),
	ticket TEXT
);

DROP TABLE IF EXISTS refresh_tokens;
CREATE TABLE refresh_tokens (
    entry_id SERIAL PRIMARY KEY,
    token bytea,
    expire_time BIGINT
);