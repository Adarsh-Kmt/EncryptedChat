CREATE TABLE user_table(
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE,
    public_key VARCHAR(255) NOT NULL
);

