CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    pseudo TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE cards (
    id SERIAL PRIMARY KEY,
    external_id TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    date TIMESTAMP NOT NULL,
    status TEXT NOT NULL,
    completed BOOLEAN NOT NULL,
    venue_name TEXT,
    city TEXT,
    region TEXT,
    country TEXT
);

CREATE TABLE fighters (
    id SERIAL PRIMARY KEY,
    external_id TEXT NOT NULL UNIQUE,
    full_name TEXT NOT NULL,
    record TEXT,
    country_name TEXT,
    country_flag TEXT,
    fighter_image TEXT
);

CREATE TABLE fights (
    id SERIAL PRIMARY KEY,
    external_id TEXT NOT NULL UNIQUE,
    card_id INTEGER NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    fighter1_id INTEGER NOT NULL REFERENCES fighters(id),
    fighter2_id INTEGER NOT NULL REFERENCES fighters(id),
    winner_fighter_id INTEGER REFERENCES fighters(id),
    category TEXT,
    status TEXT NOT NULL,
    completed BOOLEAN NOT NULL,
    points_good_prediction INTEGER NOT NULL
);

CREATE TABLE predictions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    fight_id INTEGER NOT NULL REFERENCES fights(id) ON DELETE CASCADE,
    predicted_winner_id INTEGER NOT NULL REFERENCES fighters(id),
    points_obtained INTEGER NOT NULL, --TODO: si le combat n'est pas terminé alors null semble approprié
    UNIQUE (user_id, fight_id) -- Une prédiction par user
);

CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_sessions_token_hash ON sessions(token_hash);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);