CREATE TABLE IF NOT EXISTS users
(
    id         INT UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
    first_name varchar(255),
    last_name  varchar(255),
    email      varchar(255),
    active     bool default false
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
