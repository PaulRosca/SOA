-- auth database
DROP TABLE IF EXISTS user;

CREATE TABLE user (
  id         INT AUTO_INCREMENT PRIMARY KEY,
  email      VARCHAR(255) NOT NULL,
  first_name VARCHAR(255) NOT NULL,
  last_name  VARCHAR(255) NOT NULL,
  password   VARCHAR(255) NOT NULL,
  type       ENUM('regular', 'admin') DEFAULT 'regular',
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

--orders databases
DROP TABLE IF EXISTS orders;

CREATE TABLE orders (
  id          INT AUTO_INCREMENT PRIMARY KEY,
  status      ENUM('processing', 'shipped') NOT NULL,
  address     VARCHAR(1024) NOT NULL,
  timestamp   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  email       VARCHAR(255) NOT NULL,
  user_id     INT,
  INDEX usr_ind (user_id),
  FOREIGN KEY (user_id) REFERENCES auth.user(id)
);

DROP TABLE IF EXISTS order_item;

CREATE TABLE order_item (
  id          INT AUTO_INCREMENT PRIMARY KEY,
  order_id    INT NOT NULL,
  product_id  INT NOT NULL,
  quantity    INT NOT NULL,

  INDEX ord_ind (order_id),
  INDEX prd_ind (product_id),

  FOREIGN KEY (order_id) REFERENCES orders(id),
  FOREIGN KEY (product_id) REFERENCES catalog.product(id)
);
