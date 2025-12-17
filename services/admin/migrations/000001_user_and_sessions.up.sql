CREATE TABLE users (
    id VARCHAR,
    username VARCHAR(255),
    hashed_password VARCHAR(255)
);

CREATE TABLE user_sessions (
    id          VARCHAR(255),
    user_id     VARCHAR(255),
    user_name VARCHAR(255),
    created_at  TIMESTAMP,
    expires_at  TIMESTAMP
);

INSERT INTO users(id,username,hashed_password)
VALUES ('21032d02-7757-49c6-a817-db85b10f5c4b','admin','$2a$12$YU2nW6tmRO0mlPby0gAmH.PQDzO8YMr0IJB5Ah5cdT7G59RH6Jxom')
