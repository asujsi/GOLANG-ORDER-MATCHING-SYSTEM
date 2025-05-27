This document provides a step-by-step guide to set up the Golang-based Order Matching Engine using MySQL.

---

## Prerequisites

Ensure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.21 or newer)
- [MySQL Server](https://dev.mysql.com/downloads/mysql/)

---

## STEP 1: Clone the Repository

```bash
git clone https://github.com/your-username/golang-order-matching-system.git
cd golang-order-matching-system
go mod tidy
```

---

## STEP 2: Set Up the MySQL Database

Start your MySQL server and run the following commands to create the required database and tables:

```sql
CREATE DATABASE IF NOT EXISTS testdb;
USE testdb;

CREATE TABLE IF NOT EXISTS orders (
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

CREATE TABLE IF NOT EXISTS trades (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  buy_order_id BIGINT NOT NULL,
  sell_order_id BIGINT NOT NULL,
  price DECIMAL(10,2) NOT NULL,
  quantity INT NOT NULL,
  executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## STEP 3: Configure the MySQL Connection in Code

Open the file: `internal/db/mysql.go`  
Locate the line with the DSN and update it with your local credentials:

```go
dsn := "youruser:yourpassword@tcp(127.0.0.1:3306)/testdb"
```

Replace:
- `youruser` with your MySQL username
- `yourpassword` with your MySQL password
- `testdb` with your database name if different

---

## STEP 4: Run the Application

In your terminal:

```bash
go run ./cmd
```

You should see:

```
Connected to MySQL
[GIN-debug] Listening and serving HTTP on :8080
```

---

## STEP 5: Test the API Endpoints

### Place a new buy limit order

```bash
curl -X POST http://localhost:8080/orders -H "Content-Type: application/json" -d '{"symbol":"AAPL","side":"buy","type":"limit","price":150,"quantity":10}'
```

### Cancel an order by ID

```bash
curl -X DELETE http://localhost:8080/orders/1
```

### View the current order book for a symbol

```bash
curl http://localhost:8080/orderbook?symbol=AAPL
```

---

## Verify It Works

- You receive 201 responses on valid order creation
- Orders show up in your MySQL database
- Matching engine logs trades when opposing orders exist

---

## Troubleshooting

| Issue | Possible Fix |
|-------|--------------|
| 500 Internal Server Error | Ensure tables are created and columns match |
| Ping DB failed | Check MySQL credentials and running status |
| Unknown column errors | Re-run the schema from Step 2 |
| Enum value error | Ensure you use "buy"/"sell" and "limit"/"market" exactly |

---

## Done!

You now have a working Golang Order Matching Engine with a MySQL backend, accessible over a REST API.
