-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CREATE SCHEMA IF NOT EXISTS devday;

-- PRODUCTS
CREATE TABLE IF NOT EXISTS products
(
    id          UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    create_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name        TEXT                     NOT NULL,
    kind        TEXT                     NOT NULL
);
