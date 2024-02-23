PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;

CREATE TABLE users (
    id INTEGER PRIMARY KEY ASC, 
    uuid TEXT NOT NULL UNIQUE, 
    name TEXT NOT NULL DEFAULT '',   
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL DEFAULT '',
    email_verified BOOLEAN NOT NULL DEFAULT 0,
    can_auth BOOLEAN NOT NULL DEFAULT 1,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX users_uuid ON users (uuid);
CREATE INDEX users_name ON users (name);
CREATE INDEX users_username ON users (username);
CREATE INDEX users_email ON users (email);

CREATE TRIGGER users_updated_trigger AFTER UPDATE ON users
BEGIN
      update users SET updated_on = CURRENT_TIMESTAMP WHERE id=NEW.id;
END;
