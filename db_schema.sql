CREATE TYPE bot_status AS ENUM (
    'creating',
    'aging',
    'active',
    'disabled'
);

CREATE TABLE bots (
    email varchar(255) PRIMARY KEY,
    ip inet,
    status bot_status DEFAULT 'creating',
    created_at timestamp DEFAULT NOW(),
    aging_end_date timestamp
);

