DROP SCHEMA IF EXISTS my_schema CASCADE;

CREATE SCHEMA my_schema;

SET search_path TO my_schema;

CREATE TYPE bot_status AS ENUM (
    'creating',
    'aging',
    'active',
    'disabled'
);

CREATE TABLE bots (
    email varchar(255) PRIMARY KEY,
    status bot_status DEFAULT 'creating',
    created_at timestamp DEFAULT NOW(),
    last_used timestamp DEFAULT NOW(),
    aging_end_date timestamp,
    first_name varchar(100),
    last_name varchar(100),
    username varchar(100),
    password VARCHAR(100)
);

INSERT INTO bots (email)
    VALUES ('example@gmail.com'),
    ('test@mail.me')
