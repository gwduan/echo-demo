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

INSERT INTO users(name, password, reg_date) VALUES('admin', '$argon2id$v=19$m=65536,t=1,p=12$4qYJsiDikwKPTI2p9GRxDA$os2AZJCH2X0xf6BYI0FUYOm4CZuH/kk4bew+IyM96sg', now());
