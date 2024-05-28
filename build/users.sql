PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;

CREATE TABLE roles (
    id INTEGER PRIMARY KEY ASC, 
    uuid TEXT NOT NULL UNIQUE, 
    name TEXT NOT NULL UNIQUE,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX roles_name_idx ON roles (name);

INSERT INTO roles (uuid, name) VALUES('908e6e2e-e90e-4548-a3a0-67ad356db923', '/Superuser');
INSERT INTO roles (uuid, name) VALUES('d99d8a0d-dd1c-4f9a-9736-fde7904386d8', '/Admin');
INSERT INTO roles (uuid, name) VALUES('a2296fab-a06d-4a16-a224-4f95613cf4a4', '/Standard');

CREATE TABLE permissions (
    id INTEGER PRIMARY KEY ASC, 
    uuid TEXT NOT NULL UNIQUE, 
    name TEXT NOT NULL UNIQUE,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX permissions_name_idx ON permissions (name);

INSERT INTO permissions (uuid, name) VALUES('5d224abe-bf22-4661-9ead-85cdc91746a5', '/All');

CREATE TABLE roles_permissions (
    id INTEGER PRIMARY KEY ASC, 
    role_uuid TEXT NOT NULL, 
    permission_uuid TEXT NOT NULL,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY(role_uuid) REFERENCES roles(uuid),
    FOREIGN KEY(permission_uuid) REFERENCES permissions(uuid));
CREATE INDEX roles_permissions_role_uuid_idx ON roles_permissions (role_uuid);
CREATE INDEX roles_permissions_permission_uuid_idx ON roles_permissions (permission_uuid);

INSERT INTO roles_permissions (role_uuid, permission_uuid) VALUES('908e6e2e-e90e-4548-a3a0-67ad356db923', '5d224abe-bf22-4661-9ead-85cdc91746a5');
INSERT INTO roles_permissions (role_uuid, permission_uuid) VALUES('d99d8a0d-dd1c-4f9a-9736-fde7904386d8', '5d224abe-bf22-4661-9ead-85cdc91746a5');


CREATE TABLE users (
    id INTEGER PRIMARY KEY ASC, 
    uuid TEXT NOT NULL UNIQUE, 
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL DEFAULT '',
    email_verified BOOLEAN NOT NULL DEFAULT 0,
    can_signin BOOLEAN NOT NULL DEFAULT 1,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX users_uuid ON users (uuid);
-- CREATE INDEX name ON users (first_name, last_name);
CREATE INDEX users_username ON users (username);
CREATE INDEX users_email ON users (email);

CREATE TRIGGER users_updated_trigger AFTER UPDATE ON users
BEGIN
      update users SET updated_on = CURRENT_TIMESTAMP WHERE id=NEW.id;
END;


CREATE TABLE users_roles (
    id INTEGER PRIMARY KEY ASC, 
    user_uuid TEXT NOT NULL,
    role_uuid TEXT NOT NULL, 
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY(user_uuid) REFERENCES users(uuid),
    FOREIGN KEY(role_uuid) REFERENCES roles(uuid));
CREATE INDEX users_roles_user_uuid_idx ON users_roles (user_uuid);
CREATE INDEX users_roles_role_uuid_idx ON users_roles (role_uuid);



CREATE TABLE users_sessions(
  id INTEGER PRIMARY KEY ASC,
  uuid TEXT NOT NULL,
  session_id INTEGER NOT NULL UNIQUE,
  FOREIGN KEY(uuid) REFERENCES users(uuid)
);
CREATE INDEX users_sessions_uuid ON users_sessions (uuid);
CREATE INDEX users_sessions_session_id ON users_sessions (session_id);