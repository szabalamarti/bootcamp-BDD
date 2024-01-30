-- Delete user if already exists
DROP USER IF EXISTS 'user1'@'localhost';

-- Create user with all privileges
CREATE USER 'user1'@'localhost' IDENTIFIED BY 'secret_password';
GRANT ALL PRIVILEGES ON *.* TO 'user1'@'localhost';
-- 
CREATE DATABASE  IF NOT EXISTS `supermarket_test` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `supermarket_test`;

DROP TABLE IF EXISTS `products`;
CREATE TABLE `products` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) DEFAULT NULL,
  `quantity` int DEFAULT NULL,
  `code_value` varchar(50) DEFAULT NULL,
  `is_published` varchar(50) DEFAULT NULL,
  `expiration` date DEFAULT NULL,
  `price` decimal(5,2) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `warehouses`;

CREATE TABLE `warehouses` (
  `id` int NOT NULL,
  `name` varchar(255) NOT NULL,
  `address` varchar(150) NOT NULL,
  `telephone` varchar(150) NOT NULL,
  `capacity` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Asignar la columna `id` 
ALTER TABLE `warehouses`
  ADD PRIMARY KEY (`id`);

-- Modificar tabla warehouses para que el id sea autoincrementable
ALTER TABLE `warehouses`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

-- Se agrega una columna para el id del warehouse a la tabla products
ALTER TABLE `products` ADD `id_warehouse` INT NOT NULL AFTER `price`;

-- Designar el id_warehouse con el valor 1 a cada product
UPDATE `products` SET `id_warehouse` = '1';