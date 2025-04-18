CREATE TABLE IF NOT EXISTS orders (
                                      id          BIGINT AUTO_INCREMENT PRIMARY KEY,
                                      user_id     BIGINT NOT NULL,
                                      total       DECIMAL(10,2) NOT NULL,
    status      VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
    );

CREATE TABLE IF NOT EXISTS order_items (
                                           id          BIGINT AUTO_INCREMENT PRIMARY KEY,
                                           order_id    BIGINT NOT NULL,
                                           product_id  BIGINT NOT NULL,
                                           quantity    INT NOT NULL,
                                           unit_price  DECIMAL(10,2) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id)   REFERENCES orders(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
    );
