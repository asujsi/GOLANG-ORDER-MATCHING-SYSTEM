package api

import (
	"golang-order-matching-system/internal/db"
	"golang-order-matching-system/internal/match"
	"golang-order-matching-system/internal/models"
	"log"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/orders", placeOrderHandler)
	r.DELETE("/orders/:id", cancelOrderHandler)
	r.GET("/orderbook", getOrderBookHandler)
}

// Define handler function skeletons here...
func placeOrderHandler(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	order.RemainingQuantity = order.Quantity
	order.Status = "open"

	// Insert into DB and get order ID
	res, err := db.DB.Exec(`INSERT INTO orders (symbol, side, type, price, quantity, remaining_quantity, status) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		order.Symbol, order.Side, order.Type, order.Price, order.Quantity, order.Quantity, order.Status)
	if err != nil {
		log.Println("SQL Insert Error:", err)
		c.JSON(500, gin.H{"error": "db insert error"})
		return
	}
	order.ID, _ = res.LastInsertId()

	// Run matching logic
	ob := match.GetOrCreateBook(order.Symbol)
	trades := ob.Match(&order)

	// Insert resulting trades
	for _, trade := range trades {
		_, _ = db.DB.Exec(`INSERT INTO trades (buy_order_id, sell_order_id, price, quantity) VALUES (?, ?, ?, ?)`,
			trade.BuyOrderID, trade.SellOrderID, trade.Price, trade.Quantity)
	}

	// Update order in DB
	_, _ = db.DB.Exec(`UPDATE orders SET remaining_quantity = ?, status = ? WHERE id = ?`,
		order.RemainingQuantity, order.Status, order.ID)

	// If limit order and still open, put back to book
	if order.Type == "limit" && order.RemainingQuantity > 0 {
		if order.Side == "buy" {
			ob.BuyOrders = append(ob.BuyOrders, &order)
		} else {
			ob.SellOrders = append(ob.SellOrders, &order)
		}
	}

	c.JSON(201, order)
}

func cancelOrderHandler(c *gin.Context) {
	id := c.Param("id")

	// Check order status
	var status string
	var remaining int
	err := db.DB.QueryRow("SELECT status, remaining_quantity FROM orders WHERE id = ?", id).Scan(&status, &remaining)
	if err != nil {
		c.JSON(404, gin.H{"error": "order not found"})
		return
	}

	if status == "filled" || status == "cancelled" {
		c.JSON(400, gin.H{"error": "cannot cancel a completed or already cancelled order"})
		return
	}

	// Cancel in DB
	_, err = db.DB.Exec("UPDATE orders SET status = 'cancelled', remaining_quantity = 0 WHERE id = ?", id)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to cancel order"})
		return
	}

	c.JSON(200, gin.H{"message": "order cancelled"})
}

func getOrderBookHandler(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		c.JSON(400, gin.H{"error": "symbol is required"})
		return
	}

	rows, err := db.DB.Query("SELECT id, side, price, remaining_quantity FROM orders WHERE symbol = ? AND status IN ('open', 'partially_filled') ORDER BY created_at", symbol)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch order book"})
		return
	}
	defer rows.Close()

	buyOrders := []gin.H{}
	sellOrders := []gin.H{}

	for rows.Next() {
		var id int64
		var side string
		var price float64
		var remaining int
		err := rows.Scan(&id, &side, &price, &remaining)
		if err != nil {
			continue
		}
		order := gin.H{
			"id":       id,
			"price":    price,
			"quantity": remaining,
		}
		if side == "buy" {
			buyOrders = append(buyOrders, order)
		} else {
			sellOrders = append(sellOrders, order)
		}
	}

	c.JSON(200, gin.H{
		"symbol":      symbol,
		"buy_orders":  buyOrders,
		"sell_orders": sellOrders,
	})
}
