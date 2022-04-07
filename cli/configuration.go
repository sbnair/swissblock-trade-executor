package main

import (
	"time"
)

type Configuration struct {
	OrderStream struct {
		URL string
	}

	Trader struct {
		OrderTTL time.Duration
	}

	Order struct {
		Type   string
		Symbol string
		Size   float64
		Price  float64
	}

	Database struct {
		URL string
	}
}
