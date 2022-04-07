package model

import (
	"fmt"

	"github.com/omhen/swissblock-trade-executor/v2/errors"
	"gorm.io/gorm"
)

// OrderType is the base type for enummerations to define the intention of an order
type OrderType string

func (ot OrderType) String() string {
	return string(ot)
}

const (
	Buy  OrderType = "buy"
	Sell           = "sell"
)

func ParseOrderType(value string) (ot OrderType, err error) {
	switch value {
	case "buy", "Buy", "BUY":
		ot = Buy
	case "sell", "Sell", "SELL":
		ot = Sell
	default:
		err = errors.New(fmt.Sprintf("OrderType %s not supported", value))
	}
	return
}

type Order struct {
	gorm.Model
	Type       OrderType
	Symbol     string
	Size       float64
	Pending    float64
	PriceLimit float64
	Trades     []Trade
}

func CreateOrder(db *gorm.DB, orderType OrderType, symbol string, size, priceLimit float64) (*Order, error) {
	order := Order{
		Type:       orderType,
		Symbol:     symbol,
		Size:       size,
		PriceLimit: priceLimit,
		Pending:    size,
	}
	result := db.Create(&order)
	return &order, result.Error
}
