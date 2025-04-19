CREATE TABLE IF NOT EXISTS users (
                                     id             BIGINT AUTO_INCREMENT PRIMARY KEY,
                                     username       VARCHAR(255) NOT NULL,
                                     email          VARCHAR(255) NOT NULL UNIQUE,
                                     password       VARCHAR(255) NOT NULL,
                                     first_name     VARCHAR(255)       DEFAULT NULL,
                                     last_name      VARCHAR(255)       DEFAULT NULL,
                                     country        VARCHAR(100)       DEFAULT NULL,
                                     date_of_birth  DATE               DEFAULT NULL,
                                     created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
