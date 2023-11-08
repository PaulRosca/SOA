DROP TABLE IF EXISTS user;

CREATE TABLE user (
  id         INT AUTO_INCREMENT NOT NULL,
  email      VARCHAR(255) NOT NULL,
  first_name VARCHAR(255) NOT NULL,
  last_name  VARCHAR(255) NOT NULL,
  password   VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`)
);
