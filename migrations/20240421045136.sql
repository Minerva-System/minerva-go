-- Create "companies" table
CREATE TABLE `companies` (
 `id` uuid NOT NULL DEFAULT (uuid()),
 `slug` varchar(30) NOT NULL,
 `company_name` varchar(255) NOT NULL,
 `trading_name` varchar(255) NOT NULL,
 `created_at` datetime(3) NOT NULL,
 `updated_at` datetime(3) NOT NULL,
 `deleted_at` datetime(3) NULL,
 PRIMARY KEY (`id`),
 UNIQUE INDEX `uni_companies_slug` (`slug`)
) CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
-- Create "products" table
CREATE TABLE `products` (
 `id` uuid NOT NULL DEFAULT (uuid()),
 `company_id` uuid NOT NULL,
 `description` varchar(200) NOT NULL,
 `unit` char(2) NOT NULL,
 `price` decimal(19,4) NOT NULL,
 `created_at` datetime(3) NOT NULL,
 `updated_at` datetime(3) NOT NULL,
 PRIMARY KEY (`id`, `company_id`),
 INDEX `fk_products_company` (`company_id`),
 CONSTRAINT `fk_products_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
-- Create "users" table
CREATE TABLE `users` (
 `id` uuid NOT NULL DEFAULT (uuid()),
 `company_id` uuid NOT NULL,
 `login` varchar(25) NOT NULL,
 `name` varchar(100) NOT NULL,
 `email` varchar(50) NULL,
 `pwhash` longblob NOT NULL,
 `created_at` datetime(3) NOT NULL,
 `updated_at` datetime(3) NOT NULL,
 PRIMARY KEY (`id`, `company_id`),
 INDEX `fk_users_company` (`company_id`),
 UNIQUE INDEX `uni_users_login` (`login`),
 CONSTRAINT `fk_users_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
