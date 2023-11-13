-- auth database
DROP TABLE IF EXISTS user;

CREATE TABLE user (
  id         INT AUTO_INCREMENT PRIMARY KEY,
  email      VARCHAR(255) NOT NULL,
  first_name VARCHAR(255) NOT NULL,
  last_name  VARCHAR(255) NOT NULL,
  password   VARCHAR(255) NOT NULL,
  UNIQUE KEY unique_email (email)
);

--catalog database
DROP TABLE IF EXISTS product;

CREATE TABLE product (
  id          INT AUTO_INCREMENT PRIMARY KEY,
  title       VARCHAR(128) NOT NULL,
  description VARCHAR(255) NOT NULL,
  category    VARCHAR(128) NOT NULL,
  stock       INT NOT NULL,
  price       FLOAT(7,2) NOT NULL,
  image       LONGBLOB
);
