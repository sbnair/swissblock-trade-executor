package trader

import (
	"context"
	"math"
	"time"

	"github.com/omhen/swissblock-trade-executor/v2/client"
	"github.com/omhen/swissblock-trade-executor/v2/database"
	"github.com/omhen/swissblock-trade-executor/v2/errors"
	"github.com/omhen/swissblock-trade-executor/v2/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type limitTrader struct {
	timeout time.Duration
}

func (t *limitTrader) MatchOrder(ctx context.Context, order *model.Order, reader client.BookReader) error {
	receiverChannel, errChannel := reader.StartReading(ctx, order.Symbol)
	defer reader.Stop(ctx)

	timeLimit := time.Now().Add(t.timeout)

	for {
		select {
		case err := <-errChannel:
			log.Error("While processing order: ", err)
			return errors.Annotate(err, "While processing order")
		case bookItem := <-receiverChannel:
			err := t.maybeExecuteTrade(ctx, order, bookItem)
			if err != nil {
				return errors.Annotate(err, "While attempting to trade")
			}
			if order.Pending == 0 {
				log.WithFields(log.Fields{
					"orderId": order.ID,
				}).Info("Order executed successfully")
				return nil
			}
			if time.Now().After(timeLimit) {
				log.Warn("The Order could not be matched in a timely manner")
				return nil
			}
		}
	}
}

func (t *limitTrader) maybeExecuteTrade(ctx context.Context, order *model.Order, bookItem *model.BookItem) error {
	quantity, price := t.matchedQuantityAndPrice(order, bookItem)
	if quantity == 0 || price == 0 {
		// No price nor quantity matched
		return nil
	}
	log.Infof("Match found on %s", bookItem)
	db, _ := interface{}(ctx.Value(database.ContextKeyDB)).(*gorm.DB)
	trade := model.NewTrade(quantity, price)
	order.Trades = append(order.Trades, *trade)
	order.Pending -= quantity
	db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&order)
	result := db.Save(&order)
	return result.Error
}

func (t *limitTrader) matchedQuantityAndPrice(order *model.Order, bookItem *model.BookItem) (quantity float64, price float64) {
	switch order.Type {
	case model.Buy:
		if order.PriceLimit >= bookItem.BestAskPrice {
			quantity = math.Min(order.Pending, bookItem.BestAskQuantity)
			price = bookItem.BestAskPrice
		}
	case model.Sell:
		if order.PriceLimit <= bookItem.BestBidPrice {
			quantity = math.Min(order.Pending, bookItem.BestBidQuantity)
			price = bookItem.BestBidPrice
		}
	}
	return
}

func NewLimitTrader(timeout time.Duration) Trader {
	return &limitTrader{
		timeout: timeout,
	}
}
