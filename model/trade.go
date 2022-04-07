package model

import "gorm.io/gorm"

type Trade struct {
	gorm.Model
	Quantity float64
	Price    float64
	OrderID  int
	Order    Order
}

func NewTrade(quantity, price float64) *Trade {
	return &Trade{
		Quantity: quantity,
		Price:    price,
	}
}
