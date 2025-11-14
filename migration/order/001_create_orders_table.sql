CREATE DATABASE IF NOT EXISTS edot_order CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE edot_order;

CREATE TABLE IF NOT EXISTS orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    shop_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    total_price DECIMAL(10,2) NOT NULL,
    status ENUM('pending', 'confirmed', 'shipped', 'delivered', 'cancelled') NOT NULL DEFAULT 'pending',
    order_data JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Indexes for better query performance
    INDEX idx_user_id (user_id),
    INDEX idx_shop_id (shop_id),
    INDEX idx_product_id (product_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    
    -- Constraints
    CONSTRAINT chk_quantity_positive CHECK (quantity > 0),
    CONSTRAINT chk_total_price_positive CHECK (total_price > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Add comments for documentation
ALTER TABLE orders 
COMMENT = 'Orders table for managing customer orders and order processing';

ALTER TABLE orders 
MODIFY COLUMN id INT AUTO_INCREMENT PRIMARY KEY COMMENT 'Unique order identifier',
MODIFY COLUMN user_id INT NOT NULL COMMENT 'ID of the user who placed the order',
MODIFY COLUMN shop_id INT NOT NULL COMMENT 'ID of the shop where the order was placed',
MODIFY COLUMN product_id INT NOT NULL COMMENT 'ID of the ordered product',
MODIFY COLUMN quantity INT NOT NULL DEFAULT 1 COMMENT 'Quantity of products ordered',
MODIFY COLUMN total_price DECIMAL(10,2) NOT NULL COMMENT 'Total price for the order',
MODIFY COLUMN status ENUM('pending', 'confirmed', 'shipped', 'delivered', 'cancelled') NOT NULL DEFAULT 'pending' COMMENT 'Current status of the order',
MODIFY COLUMN order_data JSON COMMENT 'Additional order data and metadata stored as JSON',
MODIFY COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Timestamp when the order was created',
MODIFY COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Timestamp when the order was last updated';