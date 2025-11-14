CREATE DATABASE IF NOT EXISTS edot_product CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE edot_product;

CREATE TABLE IF NOT EXISTS products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    stock INT NOT NULL DEFAULT 0,
    on_hold_stock INT NOT NULL DEFAULT 0,
    shop_id INT NOT NULL,
    shop_metadata JSON,
    status ENUM('active', 'inactive', 'discontinued') NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_products_name (name),
    INDEX idx_products_shop_id (shop_id),
    INDEX idx_products_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert sample products
INSERT INTO products (name, description, price, stock, on_hold_stock, shop_id, shop_metadata, status) VALUES 
('iPhone 15 Pro', 'Latest Apple iPhone with advanced camera system', 999.99, 50, 5, 1, JSON_OBJECT('shop_name', 'TechStore Pro', 'shop_id', 1, 'status', 'verified'), 'active'),
('Samsung Galaxy S24', 'Samsung flagship smartphone with AI features', 899.99, 30, 3, 1, JSON_OBJECT('shop_name', 'TechStore Pro', 'shop_id', 1, 'status', 'verified'), 'active'),
('MacBook Pro 16"', 'Apple MacBook Pro with M3 chip', 2499.99, 15, 2, 2, JSON_OBJECT('shop_name', 'Apple Premium Store', 'shop_id', 2, 'status', 'premium'), 'active'),
('Dell XPS 13', 'Ultrabook with Intel Core processor', 1299.99, 25, 0, 3, JSON_OBJECT('shop_name', 'Computer World', 'shop_id', 3, 'status', 'active'), 'active'),
('Men''s T-Shirt', 'Comfortable cotton t-shirt', 29.99, 100, 10, 4, JSON_OBJECT('shop_name', 'Fashion Hub', 'shop_id', 4, 'status', 'active'), 'active'),
('Women''s Dress', 'Elegant summer dress', 79.99, 45, 8, 4, JSON_OBJECT('shop_name', 'Fashion Hub', 'shop_id', 4, 'status', 'active'), 'active'),
('The Great Gatsby', 'Classic American novel by F. Scott Fitzgerald', 12.99, 200, 0, 5, JSON_OBJECT('shop_name', 'BookWorld', 'shop_id', 5, 'status', 'verified'), 'active'),
('Clean Code', 'A Handbook of Agile Software Craftsmanship', 45.99, 75, 5, 5, JSON_OBJECT('shop_name', 'BookWorld', 'shop_id', 5, 'status', 'verified'), 'active');

-- add table for product hold audit
CREATE TABLE IF NOT EXISTS product_hold_audit (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    status ENUM('held', 'success', 'cancelled') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_product_id (product_id),
    INDEX idx_order_id (order_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci; 