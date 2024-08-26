PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;

DROP TABLE IF EXISTS permissions;
CREATE TABLE permissions (
    id INTEGER PRIMARY KEY ASC, 
    public_id TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT "",
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX roles_name_idx ON permissions (name);

INSERT INTO permissions (public_id, name, description) VALUES('uwkrk2ljj387', 'super', 'Superuser');
INSERT INTO permissions (public_id, name, description) VALUES('iz4kbfy3z0a3', 'admin', 'Administrator');
INSERT INTO permissions (public_id, name, description) VALUES('loq75e7zqcbl', 'user', 'User');
INSERT INTO permissions (public_id, name, description) VALUES('kflynb03pxbj', 'login', 'Can login');
INSERT INTO permissions (public_id, name, description) VALUES('og1o5d0p0mjy', 'rdf', 'Can view RDF lab data');

DROP TABLE IF EXISTS roles;
CREATE TABLE roles (
    id INTEGER PRIMARY KEY ASC, 
    public_id TEXT NOT NULL UNIQUE, 
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT "",
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX permissions_name_idx ON permissions (name);

INSERT INTO roles (public_id, name) VALUES('p1gbjods0h90', 'super');
INSERT INTO roles (public_id, name) VALUES('mk4bgg4w43fp', 'admin');
INSERT INTO roles (public_id, name) VALUES('3xvte0ik4aq4', 'user');
-- INSERT INTO roles (public_id, name) VALUES('UZuAVHDGToa4F786IPTijA==', 'GetDNA');
INSERT INTO roles (public_id, name) VALUES('x4ewk9papip2', 'login');
INSERT INTO roles (public_id, name) VALUES('kh2yynyheqhv', 'rdf');

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
    public_id TEXT NOT NULL UNIQUE, 
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL DEFAULT '',
    email_verified BOOLEAN NOT NULL DEFAULT 0,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX users_uuid ON users (public_id);
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
  public_id TEXT NOT NULL,
  session_id INTEGER NOT NULL UNIQUE,
  FOREIGN KEY(public_id) REFERENCES users(public_id)
);
CREATE INDEX users_sessions_uuid ON users_sessions (public_id);
CREATE INDEX users_sessions_session_id ON users_sessions (session_id);