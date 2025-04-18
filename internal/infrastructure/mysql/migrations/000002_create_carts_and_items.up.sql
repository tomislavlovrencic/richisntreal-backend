CREATE TABLE IF NOT EXISTS carts (
                                     id         BIGINT AUTO_INCREMENT PRIMARY KEY,
                                     user_id    BIGINT NOT NULL,
                                     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS cart_items (
                                          id          BIGINT AUTO_INCREMENT PRIMARY KEY,
                                          cart_id     BIGINT NOT NULL,
                                          product_id  BIGINT NOT NULL,
                                          quantity    INT NOT NULL DEFAULT 1,
                                          unit_price  DECIMAL(10,2) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (cart_id) REFERENCES carts(id) ON DELETE CASCADE
    );
