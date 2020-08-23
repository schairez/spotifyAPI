CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    provider_type provider_type NOT NULL DEFAULT 'SPOTIFY',
    user_provider_id TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    image_url TEXT NOT NULL UNIQUE,
    num_followers INT NOT NULL DEFAULT 0,
    country VARCHAR(2) NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
CREATE TYPE provider_type AS ENUM ('SPOTIFY', 'YOUTUBE');
CREATE TYPE album_type AS ENUM ('ALBUM', 'SINGLE');
CREATE TABLE IF NOT EXISTS likedsongs (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    song_id INT NOT NULL,
    added_at timestamptz NOT NULL,
    popularity INT NOT NULL,
    preview_url TEXT NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(song_id) REFERENCES songs(id),
    UNIQUE (user_id, song_id)
);
CREATE TABLE IF NOT EXISTS songs(
    id SERIAL PRIMARY KEY,
    song_name TEXT,
    popularity INT CHECK (
        testing > 0
        AND testing < 100
    )
);
CREATE TABLE IF NOT EXISTS albums(
    id SERIAL PRIMARY KEY,
    album_name TEXT NOT NULL,
    album_type album_type NOT NULL DEFAULT 'ALBUM'
);
CREATE TABLE IF NOT EXISTS artists(
    id SERIAL PRIMARY KEY,
    artist TEXT NOT NULL
);