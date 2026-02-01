CREATE DATABASE IF NOT EXISTS pos_test
CHARACTER SET utf8mb4
COLLATE utf8mb4_general_ci;

USE pos_test;

DROP TABLE IF EXISTS products;
CREATE TABLE products (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    stock_qty INT NOT NULL
) ENGINE=InnoDB;

INSERT INTO products (name, price, stock_qty)
VALUES
("Shoes", 1150, 20),
("Belt", 450, 3),
("T-Shirt", 950, 76);

DROP TABLE IF EXISTS branches;
CREATE TABLE branches (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
) ENGINE=InnoDB;

INSERT INTO branches (name, is_active) 
VALUES
("B-Test-1", TRUE), 
("B-Test-2", TRUE);

DROP TABLE IF EXISTS store_products;
CREATE TABLE store_products (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    branch_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    qty INT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT fk_store_products_branch
        FOREIGN KEY (branch_id)
        REFERENCES branches(id),

    CONSTRAINT fk_store_product
        FOREIGN KEY (product_id)
        REFERENCES products(id)
) ENGINE=InnoDB;

INSERT INTO store_products (branch_id, product_id, qty, price, is_active) 
VALUES
(1, 1, 546, 1150, TRUE),
(1, 2, 870, 450, TRUE), 
(1, 3, 209, 950, FALSE),
(2, 1, 146, 1150, FALSE),
(2, 2, 1070, 450, TRUE), 
(2, 3, 909, 950, TRUE);

DROP TABLE IF EXISTS orders;
CREATE TABLE orders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    branch_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    qty INT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    total DECIMAL(10,2) NOT NULL,

    CONSTRAINT fk_orders_branch
        FOREIGN KEY (branch_id)
        REFERENCES branches(id),

    CONSTRAINT fk_orders_product
        FOREIGN KEY (product_id)
        REFERENCES products(id)
) ENGINE=InnoDB;

INSERT INTO orders (branch_id, product_id, qty, price, total)
VALUES
(1, 1, 2, 1150, 2300),
(1, 2, 1, 450, 450),
(1, 3, 3, 950, 2850),
(1, 1, 1, 1150, 1150),
(1, 3, 1, 950, 950),
(2, 2, 2, 450, 900),
(2, 3, 1, 950, 950),
(2, 1, 1, 1150, 1150);
