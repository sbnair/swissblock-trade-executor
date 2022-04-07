package reports

import (
	"github.com/omhen/swissblock-trade-executor/v2/model"
	log "github.com/sirupsen/logrus"
)

func OrderLogSummary(order *model.Order) {
	if len(order.Trades) == 0 {
		log.WithFields(log.Fields{
			"orderId":    order.ID,
			"type":       order.Type,
			"symbol":     order.Symbol,
			"size":       order.Size,
			"priceLimit": order.PriceLimit,
		}).Info("The order did not generate any trade")
		return
	}
	log.WithFields(log.Fields{
		"orderId":    order.ID,
		"type":       order.Type,
		"symbol":     order.Symbol,
		"size":       order.Size,
		"priceLimit": order.PriceLimit,
	}).Infof("The order generated %d trades", len(order.Trades))
	for _, trade := range order.Trades {
		log.Infof(
			"id: %d, quantity: %f, price: %f, timestamp: %v",
			trade.ID,
			trade.Quantity,
			trade.Price,
			trade.CreatedAt,
		)
	}
}
