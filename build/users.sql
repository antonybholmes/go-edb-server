PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;

DROP TABLE IF EXISTS roles;
CREATE TABLE roles (
    id INTEGER PRIMARY KEY ASC, 
    uuid TEXT NOT NULL UNIQUE, 
    name TEXT NOT NULL UNIQUE,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(uuid, name));
CREATE INDEX roles_name_idx ON roles (name);

INSERT INTO roles (uuid, name) VALUES('908e6e2e-e90e-4548-a3a0-67ad356db923', 'Superuser');
INSERT INTO roles (uuid, name) VALUES('d99d8a0d-dd1c-4f9a-9736-fde7904386d8', 'Admin');
INSERT INTO roles (uuid, name) VALUES('a2296fab-a06d-4a16-a224-4f95613cf4a4', 'Standard');
INSERT INTO roles (uuid, name) VALUES('0f988db6-810f-4e0a-82f5-2493baf6b49e', 'Mutations');

DROP TABLE IF EXISTS permissions;
CREATE TABLE permissions (
    id INTEGER PRIMARY KEY ASC, 
    uuid TEXT NOT NULL, 
    name TEXT NOT NULL,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(uuid, name));
CREATE INDEX permissions_name_idx ON permissions (name);

INSERT INTO permissions (uuid, name) VALUES('5d224abe-bf22-4661-9ead-85cdc91746a5', 'Admin');
INSERT INTO permissions (uuid, name) VALUES('4a0730a9-211f-48b9-bb65-803abeca9e31', 'GetDNA');
INSERT INTO permissions (uuid, name) VALUES('7df054ba-ef7b-4240-9b40-ff537904990b', 'GetMutations');

DROP TABLE IF EXISTS role_permissions;
CREATE TABLE role_permissions (
    id INTEGER PRIMARY KEY ASC, 
    role_uuid TEXT NOT NULL, 
    permission_uuid TEXT NOT NULL,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(role_uuid, permission_uuid),
    FOREIGN KEY(role_uuid) REFERENCES roles(uuid),
    FOREIGN KEY(permission_uuid) REFERENCES permissions(uuid));
CREATE INDEX role_permissions_role_uuid_idx ON role_permissions (role_uuid);
CREATE INDEX role_permissions_permission_uuid_idx ON role_permissions (permission_uuid);

-- super/user admin
INSERT INTO role_permissions (role_uuid, permission_uuid) VALUES('908e6e2e-e90e-4548-a3a0-67ad356db923', '5d224abe-bf22-4661-9ead-85cdc91746a5');
INSERT INTO role_permissions (role_uuid, permission_uuid) VALUES('d99d8a0d-dd1c-4f9a-9736-fde7904386d8', '5d224abe-bf22-4661-9ead-85cdc91746a5');

--
-- standard
--
-- dna
INSERT INTO role_permissions (role_uuid, permission_uuid) VALUES('a2296fab-a06d-4a16-a224-4f95613cf4a4', '4a0730a9-211f-48b9-bb65-803abeca9e31');

-- mutations
INSERT INTO role_permissions (role_uuid, permission_uuid) VALUES('0f988db6-810f-4e0a-82f5-2493baf6b49e', '7df054ba-ef7b-4240-9b40-ff537904990b');

DROP TABLE IF EXISTS users;
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

DROP TABLE IF EXISTS user_roles;
CREATE TABLE user_roles (
    id INTEGER PRIMARY KEY ASC, 
    user_uuid TEXT NOT NULL,
    role_uuid TEXT NOT NULL, 
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(user_uuid, role_uuid),
    FOREIGN KEY(user_uuid) REFERENCES users(uuid),
    FOREIGN KEY(role_uuid) REFERENCES roles(uuid));
CREATE INDEX user_roles_user_uuid_idx ON user_roles (user_uuid);
CREATE INDEX user_roles_role_uuid_idx ON user_roles (role_uuid);



CREATE TABLE users_sessions(
  id INTEGER PRIMARY KEY ASC,
  uuid TEXT NOT NULL,
  session_id INTEGER NOT NULL UNIQUE,
  FOREIGN KEY(uuid) REFERENCES users(uuid)
);
CREATE INDEX users_sessions_uuid ON users_sessions (uuid);
CREATE INDEX users_sessions_session_id ON users_sessions (session_id);