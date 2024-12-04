CREATE TABLE users (
  id          SERIAL PRIMARY KEY,
  username    VARCHAR(25) NOT NULL UNIQUE,
  phone       VARCHAR(10) NOT NULL UNIQUE,
  password    VARCHAR(256) NOT NULL,
  photo       VARCHAR(512),
  description VARCHAR(256),
  role_id     INTEGER NOT NULL DEFAULT 1,
  created_at  TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE roles (
  id          SERIAL PRIMARY KEY,
  name        VARCHAR (50) NOT NULL,
  permissions TEXT[] NOT NULL DEFAULT '{}'
);

INSERT INTO roles (name, permissions)
VALUES ('Қолданушы', '{}');

-- codes

CREATE TABLE codes (
  code       VARCHAR(4) NOT NULL,
  phone      VARCHAR(10) NOT NULL,
  ip         TEXT NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- sessions

CREATE TABLE sessions (
  access_token TEXT NOT NULL UNIQUE,
  user_id INTEGER NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);