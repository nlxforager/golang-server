

CREATE TABLE signatures (
    id SERIAL PRIMARY KEY,
    key  VARCHAR NOT NULL UNIQUE,
    signature VARCHAR NOT NULL
);


CREATE TABLE accounts
(
    id              SERIAL PRIMARY KEY,
    username        VARCHAR(50) NOT NULL UNIQUE,
    hashed_password VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE emails
(
    id         SERIAL PRIMARY KEY,
    account_id VARCHAR(50) REFERENCES accounts (id),
    email      VARCHAR(100) NOT NULL UNIQUE,
);

--
-- CREATE TABLE mobiles
-- (
--     id     SERIAL PRIMARY KEY,
--     number VARCHAR REFERENCES accounts (id)
-- );
