package main

import (
	"time"

	"github.com/urfave/cli"
)

func newApp() (commandLine *cli.App) {
	commandLine = cli.NewApp()
	commandLine.Name = "Trade executor for Swissblock orders"

	conf := new(Configuration)

	commandLine.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "order_book_url",
			Usage:       "The websocket URL where to read the order book stream",
			EnvVar:      "ORDER_BOOK_URL",
			Destination: &conf.OrderStream.URL,
			Value:       "wss://stream.binance.com:9443/ws",
		},
		cli.DurationFlag{
			Name:        "order_time_limit",
			Usage:       "Max time to try to execute the order. Defaults (3600s)",
			EnvVar:      "ORDER_TIME_LIMIT",
			Destination: &conf.Trader.OrderTTL,
			Value:       time.Second * 3600,
		},
		cli.StringFlag{
			Name:        "symbol",
			Usage:       "The symbol pair to be traded.",
			EnvVar:      "ORDER_SYMBOL",
			Destination: &conf.Order.Symbol,
		},
		cli.StringFlag{
			Name:        "order_type",
			Usage:       "The type of the order. Can be either 'buy' or 'sell'. Defaults to 'buy'",
			EnvVar:      "ORDER_TYPE",
			Destination: &conf.Order.Type,
			Value:       "buy",
		},
		cli.Float64Flag{
			Name:        "size",
			Usage:       "The amount to be traded.",
			EnvVar:      "ORDER_SIZE",
			Destination: &conf.Order.Size,
		},
		cli.Float64Flag{
			Name:        "price",
			Usage:       "The price limit for the trade.",
			EnvVar:      "ORDER_PRICE_LIMIT",
			Destination: &conf.Order.Price,
		},
		cli.StringFlag{
			Name:        "database_url",
			Usage:       "The database where to store the operations and their related transactions",
			EnvVar:      "DATABASE_URL",
			Destination: &conf.Database.URL,
		},
	}

	commandLine.Action = func(c *cli.Context) {
		start(conf)
	}

	return
}
