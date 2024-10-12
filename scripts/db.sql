CREATE DATABASE echo_demo;
USE echo_demo;

CREATE TABLE users (
  id		BIGINT AUTO_INCREMENT NOT NULL,
  name		VARCHAR(128) NOT NULL,
  password	VARCHAR(255) NOT NULL,
  age		INT,
  reg_date	DATETIME NOT NULL,
  PRIMARY KEY(`id`),
  UNIQUE(`name`)
);

INSERT INTO users(name, password, reg_date) VALUES('admin', 'admin', now());
