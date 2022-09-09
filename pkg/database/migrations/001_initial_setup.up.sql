CREATE SCHEMA IF NOT EXISTS devday;

-- USERS
CREATE TYPE role AS ENUM ('guest', 'member', 'admin');

CREATE TABLE IF NOT EXISTS users
(
    id          UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    role        role                     NOT NULL DEFAULT 'guest',
    create_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name        TEXT                     NOT NULL
);


-- PRODUCTS
CREATE TABLE IF NOT EXISTS products
(
    id          UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    create_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name        TEXT                     NOT NULL
);
