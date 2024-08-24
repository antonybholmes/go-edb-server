PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;

DROP TABLE IF EXISTS roles;
CREATE TABLE roles (
    id INTEGER PRIMARY KEY ASC, 
    uuid TEXT NOT NULL UNIQUE,
    publicId TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(id, name));
CREATE INDEX roles_name_idx ON roles (name);

INSERT INTO roles (uuid, publicId, name) VALUES('908e6e2e-e90e-4548-a3a0-67ad356db923', 'su', 'Superuser');
INSERT INTO roles (uuid, publicId, name) VALUES('d99d8a0d-dd1c-4f9a-9736-fde7904386d8', 'admin', 'Administrator');
INSERT INTO roles (uuid, publicId, name) VALUES('a2296fab-a06d-4a16-a224-4f95613cf4a4', 'user', 'User');
INSERT INTO roles (uuid, publicId, name) VALUES('0f988db6-810f-4e0a-82f5-2493baf6b49e', 'mutations', 'Mutations');

DROP TABLE IF EXISTS groups;
CREATE TABLE groups (
    id INTEGER PRIMARY KEY ASC, 
    uuid TEXT NOT NULL UNIQUE, 
    name TEXT NOT NULL UNIQUE,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(id, name));
CREATE INDEX permissions_name_idx ON permissions (name);

INSERT INTO groups (uuid, name) VALUES('5d224abe-bf22-4661-9ead-85cdc91746a5', 'Superuser');
INSERT INTO groups (uuid, name) VALUES('5085ff97-1773-4496-bed5-097d2ca48ac6', 'Administrator');
INSERT INTO groups (uuid, name) VALUES('286ee0cb-ba0d-4442-ad18-05f585c2b257', 'User');
-- INSERT INTO groups (uuid, name) VALUES('4a0730a9-211f-48b9-bb65-803abeca9e31', 'GetDNA');
INSERT INTO groups (uuid, name) VALUES('7df054ba-ef7b-4240-9b40-ff537904990b', 'Mutations');

DROP TABLE IF EXISTS groups_roles;
CREATE TABLE groups_roles (
    id INTEGER PRIMARY KEY ASC, 
    group_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(group_id, role_id),
    FOREIGN KEY(group_id) REFERENCES groups(id),
    FOREIGN KEY(role_id) REFERENCES roles(id));
CREATE INDEX groups_roles_group_id_idx ON groups_roles (group_id, role_id);

-- super/user admin
INSERT INTO groups_roles (group_id, role_id) VALUES(1, 1);
INSERT INTO groups_roles (group_id, role_id) VALUES(1, 2);
INSERT INTO groups_roles (group_id, role_id) VALUES(2, 2);

--
-- standard
INSERT INTO groups_roles (group_id, role_id) VALUES(3, 3);

-- mutations
INSERT INTO groups_roles (role_uuid, permission_uuid) VALUES(4, 4);

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

DROP TABLE IF EXISTS users_groups;
CREATE TABLE users_groups (
    id INTEGER PRIMARY KEY ASC, 
    user_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL, 
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(user_id, group_id),
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(group_id) REFERENCES groups(id));
CREATE INDEX users_groups_user_id_idx ON users_groups (user_id, group_id);
 



CREATE TABLE users_sessions(
  id INTEGER PRIMARY KEY ASC,
  uuid TEXT NOT NULL,
  session_id INTEGER NOT NULL UNIQUE,
  FOREIGN KEY(uuid) REFERENCES users(uuid)
);
CREATE INDEX users_sessions_uuid ON users_sessions (uuid);
CREATE INDEX users_sessions_session_id ON users_sessions (session_id);