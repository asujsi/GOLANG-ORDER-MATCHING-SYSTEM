package match

import (
	"golang-order-matching-system/internal/models"
	"sort"
)

var OrderBooks = make(map[string]*OrderBook)

type OrderBook struct {
	BuyOrders  []*models.Order
	SellOrders []*models.Order
}

func GetOrCreateBook(symbol string) *OrderBook {
	if ob, ok := OrderBooks[symbol]; ok {
		return ob
	}
	OrderBooks[symbol] = &OrderBook{}
	return OrderBooks[symbol]
}

func (ob *OrderBook) Match(order *models.Order) []*models.Trade {
	var trades []*models.Trade

	if order.Side == "buy" {
		sort.SliceStable(ob.SellOrders, func(i, j int) bool {
			return ob.SellOrders[i].Price < ob.SellOrders[j].Price
		})

		for len(ob.SellOrders) > 0 && order.RemainingQuantity > 0 {
			bestAsk := ob.SellOrders[0]
			match := false

			if order.Type == "market" || order.Price >= bestAsk.Price {
				matchQty := min(order.RemainingQuantity, bestAsk.RemainingQuantity)
				tradePrice := bestAsk.Price
				trades = append(trades, &models.Trade{
					BuyOrderID:  order.ID,
					SellOrderID: bestAsk.ID,
					Quantity:    matchQty,
					Price:       tradePrice,
				})

				order.RemainingQuantity -= matchQty
				bestAsk.RemainingQuantity -= matchQty
				updateStatus(order)
				updateStatus(bestAsk)

				if bestAsk.RemainingQuantity == 0 {
					ob.SellOrders = ob.SellOrders[1:]
				} else {
					break
				}
				match = true
			}

			if !match {
				break
			}
		}

	} else { // sell order
		sort.SliceStable(ob.BuyOrders, func(i, j int) bool {
			return ob.BuyOrders[i].Price > ob.BuyOrders[j].Price
		})

		for len(ob.BuyOrders) > 0 && order.RemainingQuantity > 0 {
			bestBid := ob.BuyOrders[0]
			match := false

			if order.Type == "market" || order.Price <= bestBid.Price {
				matchQty := min(order.RemainingQuantity, bestBid.RemainingQuantity)
				tradePrice := bestBid.Price
				trades = append(trades, &models.Trade{
					BuyOrderID:  bestBid.ID,
					SellOrderID: order.ID,
					Quantity:    matchQty,
					Price:       tradePrice,
				})

				order.RemainingQuantity -= matchQty
				bestBid.RemainingQuantity -= matchQty
				updateStatus(order)
				updateStatus(bestBid)

				if bestBid.RemainingQuantity == 0 {
					ob.BuyOrders = ob.BuyOrders[1:]
				} else {
					break
				}
				match = true
			}

			if !match {
				break
			}
		}
	}

	return trades
}

func updateStatus(order *models.Order) {
	if order.RemainingQuantity == 0 {
		order.Status = "filled"
	} else if order.RemainingQuantity < order.Quantity {
		order.Status = "partially_filled"
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
