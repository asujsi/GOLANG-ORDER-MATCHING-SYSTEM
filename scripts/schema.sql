-- orders table
CREATE TABLE orders (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  symbol VARCHAR(20) NOT NULL,
  side ENUM('buy', 'sell') NOT NULL,
  type ENUM('limit', 'market') NOT NULL,
  price DECIMAL(10,2),
  quantity INT NOT NULL,
  remaining_quantity INT NOT NULL,
  status ENUM('open', 'partially_filled', 'filled', 'cancelled') DEFAULT 'open',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- trades table
CREATE TABLE trades (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  buy_order_id BIGINT NOT NULL,
  sell_order_id BIGINT NOT NULL,
  price DECIMAL(10,2) NOT NULL,
  quantity INT NOT NULL,
  executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
