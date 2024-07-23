CREATE TABLE strong_password_log (
    id BIGSERIAL PRIMARY KEY,
    req JSONB NOT NULL,
    res JSONB NOT NULL
);
