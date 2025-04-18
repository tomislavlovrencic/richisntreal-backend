CREATE TABLE IF NOT EXISTS payment_transactions (
                                                    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
                                                    order_id        BIGINT NOT NULL,
                                                    amount          DECIMAL(10,2) NOT NULL,
    currency        VARCHAR(3) NOT NULL DEFAULT 'USD',
    provider        VARCHAR(50) NOT NULL,         -- e.g. "stripe", "paypal"
    provider_tx_id  VARCHAR(255) DEFAULT NULL,    -- gateway’s transaction ID
    token           VARCHAR(255) NOT NULL,        -- card token or payment method ID
    status          VARCHAR(50) NOT NULL,         -- e.g. "pending", "succeeded", "failed"
    failure_message TEXT,                         -- populated on error
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id)
    );
