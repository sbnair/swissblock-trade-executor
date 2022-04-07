package model

import "fmt"

type BookItem struct {
	UpdateID        uint    `json:"u,omitempty"`
	Symbol          string  `json:"s,omitempty"`
	BestBidPrice    float64 `json:"b,string,omitempty"`
	BestBidQuantity float64 `json:"B,string,omitempty"`
	BestAskPrice    float64 `json:"a,string,omitempty"`
	BestAskQuantity float64 `json:"A,string,omitempty"`
}

func (bi *BookItem) String() string {
	return fmt.Sprintf(
		"BookItem<uid: %d, symbol: %s, b: %f, B: %f, a: %f, A: %f>",
		bi.UpdateID,
		bi.Symbol,
		bi.BestBidPrice,
		bi.BestBidQuantity,
		bi.BestAskPrice,
		bi.BestAskQuantity,
	)
}
