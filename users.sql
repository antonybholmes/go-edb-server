PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;

DROP TABLE IF EXISTS permissions;
CREATE TABLE permissions (
    id INTEGER PRIMARY KEY ASC, 
    uuid TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT "",
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX roles_name_idx ON permissions (name);

INSERT INTO permissions (uuid, name, description) VALUES('uwkrk2ljj387', 'SU', 'Superuser');
INSERT INTO permissions (uuid, name, description) VALUES('iz4kbfy3z0a3', 'Admin', 'Administrator');
INSERT INTO permissions (uuid, name, description) VALUES('loq75e7zqcbl', 'User', 'User');
INSERT INTO permissions (uuid, name, description) VALUES('kflynb03pxbj', 'Login', 'Can login');
INSERT INTO permissions (uuid, name, description) VALUES('og1o5d0p0mjy', 'RDF', 'Can view RDF lab data');

DROP TABLE IF EXISTS roles;
CREATE TABLE roles (
    id INTEGER PRIMARY KEY ASC, 
    uuid TEXT NOT NULL UNIQUE, 
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT "",
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX permissions_name_idx ON permissions (name);

INSERT INTO roles (uuid, name) VALUES('p1gbjods0h90', 'Superuser');
INSERT INTO roles (uuid, name) VALUES('mk4bgg4w43fp', 'Administrator');
INSERT INTO roles (uuid, name) VALUES('3xvte0ik4aq4', 'User');
-- INSERT INTO roles (uuid, name) VALUES('UZuAVHDGToa4F786IPTijA==', 'GetDNA');
INSERT INTO roles (uuid, name) VALUES('x4ewk9papip2', 'Login');
INSERT INTO roles (uuid, name) VALUES('kh2yynyheqhv', 'RDF');

DROP TABLE IF EXISTS roles_permissions;
CREATE TABLE roles_permissions (
    id INTEGER PRIMARY KEY ASC, 
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(role_id, permission_id),
    FOREIGN KEY(role_id) REFERENCES roles(id),
    FOREIGN KEY(permission_id) REFERENCES permissions(id));
CREATE INDEX roles_permissions_role_id_idx ON roles_permissions (role_id, permission_id);

-- super/user admin
INSERT INTO roles_permissions (role_id, permission_id) VALUES(1, 1);
INSERT INTO roles_permissions (role_id, permission_id) VALUES(1, 2);
INSERT INTO roles_permissions (role_id, permission_id) VALUES(2, 2);

--
-- standard
INSERT INTO roles_permissions (role_id, permission_id) VALUES(3, 3);

-- users can login
INSERT INTO roles_permissions (role_id, permission_id) VALUES(4, 4);

-- rdf
INSERT INTO roles_permissions (role_id, permission_id) VALUES(5, 5);

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

DROP TABLE IF EXISTS users_roles;
CREATE TABLE users_roles (
    id INTEGER PRIMARY KEY ASC, 
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL, 
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(user_id, role_id),
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(role_id) REFERENCES roles(id));
CREATE INDEX users_roles_user_id_idx ON users_roles (user_id, role_id);
 



CREATE TABLE users_sessions(
  id INTEGER PRIMARY KEY ASC,
  uuid TEXT NOT NULL,
  session_id INTEGER NOT NULL UNIQUE,
  FOREIGN KEY(uuid) REFERENCES users(uuid)
);
CREATE INDEX users_sessions_uuid ON users_sessions (uuid);
CREATE INDEX users_sessions_session_id ON users_sessions (session_id);