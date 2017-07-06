#SET SESSION time_zone = "+0:00";

use gbs;

DROP TABLE IF EXISTS project_team;
DROP TABLE IF EXISTS project;
DROP TABLE IF EXISTS account_invoices;
DROP TABLE IF EXISTS invoice;
DROP TABLE IF EXISTS account;
DROP TABLE IF EXISTS customer_addresses;
DROP TABLE IF EXISTS address;
DROP TABLE IF EXISTS ref_address_type;
DROP TABLE IF EXISTS user_permission;
DROP TABLE IF EXISTS permission;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS team;
DROP TABLE IF EXISTS user;

CREATE TABLE user (
  `user_id` bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `email` VARCHAR(255) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `first_name` VARCHAR(100) NOT NULL,
  `last_name` VARCHAR(100) NOT NULL,
  `github_name` VARCHAR(255),
  `slack_name` VARCHAR(255),
  `date_became_customer` DATE,
  INDEX(user_id),
  INDEX(email),
   UNIQUE KEY (email)

) ENGINE=INNODB;



CREATE TABLE `team` (
  `team_id` bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `team_name` VARCHAR(255),
  `team_desc` VARCHAR(255),
  `created_by` bigint unsigned,
    INDEX(team_id)
   ) ENGINE=INNODB;

CREATE TABLE `team_members` (
    `teams_members_id` bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `user_id` bigint unsigned,
    `team_id` bigint unsigned,
	 INDEX(user_id),
	INDEX(team_id),
	INDEX(teams_members_id),

      CONSTRAINT `fk_team_member` FOREIGN KEY (user_id) REFERENCES user (user_id) ON UPDATE CASCADE,
      CONSTRAINT `fk_team` FOREIGN KEY (team_id) REFERENCES team (team_id) ON UPDATE CASCADE
) ENGINE=INNODB;


CREATE TABLE `permission` (
  `permission_id` bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `permissions` VARCHAR(20),
    `permission_name` VARCHAR(100),
    INDEX(permission_id)
) ENGINE=INNODB;

CREATE TABLE `user_permission` (

   `user_permission_id` bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `user_id` bigint unsigned,
  `permission_id` bigint unsigned,

	INDEX(user_permission_id)
) ENGINE=INNODB;

CREATE TABLE `ref_address_type` (
  `address_type_code` bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(100),
  `address_description` VARCHAR(30),
  	INDEX (address_type_code)
) ENGINE=INNODB;

CREATE TABLE `address` (
  `address_id` bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `user_id` bigint unsigned,
  `line_1` VARCHAR(255),
  `line_2` VARCHAR(255),
  `line_3` VARCHAR(255),
  `line_4` VARCHAR(255),
  `city` VARCHAR(255),
  `zip_or_post` INT,
  `country` VARCHAR(255),
  `state` VARCHAR(255),
  INDEX (address_id),
    INDEX (user_id),
     INDEX (country),
      INDEX (state),
   CONSTRAINT `fk_user_id` FOREIGN KEY (user_id) REFERENCES user (user_id) ON UPDATE CASCADE
) ENGINE=INNODB;

CREATE TABLE `customer_addresses` (
  `customer_address_id` bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `address_id` bigint unsigned,
  `address_type_code` bigint unsigned,
  	INDEX (customer_address_id),
	INDEX (address_id),
   CONSTRAINT `fk_user_address` FOREIGN KEY (address_id) REFERENCES address (address_id) ON UPDATE CASCADE,
   CONSTRAINT `fk_address_type`FOREIGN KEY (address_type_code) REFERENCES ref_address_type (address_type_code) ON UPDATE CASCADE
) ENGINE=INNODB;

CREATE TABLE `account` (
  `account_id` bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `user_id` bigint unsigned,
  `billing_address` bigint unsigned,
  `shipping_address` bigint unsigned,
	INDEX (account_id),
	INDEX (user_id),
	CONSTRAINT `fk_user_account` FOREIGN KEY (user_id) REFERENCES user (user_id) ON UPDATE CASCADE ,
	CONSTRAINT `fk_user_billing_address` FOREIGN KEY (billing_address) REFERENCES customer_addresses (customer_address_id) ,
	CONSTRAINT `fk_user_shipping_address` FOREIGN KEY (shipping_address) REFERENCES customer_addresses (customer_address_id)
) ENGINE=INNODB;

CREATE TABLE `invoice` (
  `invoice_id` bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `invoice_date` DATE,
  `due_date` DATE,
  `pay_date` DATE,
  `units` INT,
  `unit_price` INT,
  `description` BLOB,
  `amount_due` INT,
  `payment_amount` INT,
  `notes` BLOB,
  `next_bill_date` DATE,
   INDEX (invoice_id),
   INDEX (invoice_date),
  INDEX (due_date)
) ENGINE=INNODB;


CREATE TABLE `account_invoices`  (
   `account_invoices_id` bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `account_id` bigint unsigned,
  `invoice_id` bigint unsigned,
   INDEX (account_id),
    INDEX (invoice_id),
    INDEX (account_invoices_id)
) ENGINE=INNODB;

CREATE TABLE `project` (
  `project_id` bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(100),
  `description` BLOB,
  `difficulty` int,
    INDEX (project_id)
) ENGINE=INNODB;


CREATE TABLE `project_team` (
   `project_team_id` bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `project_id` bigint unsigned,
  `team_id` bigint unsigned,
  `project_url` VARCHAR(255),
   INDEX (project_id),
    INDEX (team_id),
     INDEX (project_team_id)
) ENGINE=INNODB;






