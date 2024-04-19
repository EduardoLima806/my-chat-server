SELECT 'CREATE DATABASE chat_server'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'chat_server')\gexec

\connect chat_server

CREATE TABLE IF NOT EXISTS app_user (
  id serial,
  username varchar(50) NOT NULL,
  displayname varchar(255),
  email varchar(150) NOT NULL,
  password varchar(150) NOT NULL,
  created timestamp NOT NULL,
  PRIMARY KEY (id)
)\gexec