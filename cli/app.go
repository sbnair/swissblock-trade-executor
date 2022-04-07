package main

import (
	"context"
	"os"

	"github.com/omhen/swissblock-trade-executor/v2/client"
	"github.com/omhen/swissblock-trade-executor/v2/database"
	"github.com/omhen/swissblock-trade-executor/v2/model"
	"github.com/omhen/swissblock-trade-executor/v2/reports"
	"github.com/omhen/swissblock-trade-executor/v2/trader"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func start(conf *Configuration) {
	db, err := database.NewDBConnection(conf.Database.URL)
	if err != nil {
		log.Panic("While connecting to database. ", err)
		os.Exit(1)
	}

	ctx := context.WithValue(context.Background(), database.ContextKeyDB, db)

	bookReader := client.NewBookReaderBinance(conf.OrderStream.URL)
	order, err := createOrder(conf, db)
	if err != nil {
		log.Panic("While saving order in DB. ", err)
		os.Exit(1)
	}

	trader := trader.NewLimitTrader(conf.Trader.OrderTTL)

	err = trader.MatchOrder(ctx, order, bookReader)
	if err != nil {
		log.Error("While trying to match the order: ", err)
		os.Exit(1)
	}

	reports.OrderLogSummary(order)
	log.Info("Process completed.")
}

func createOrder(conf *Configuration, db *gorm.DB) (*model.Order, error) {
	orderType, err := model.ParseOrderType(conf.Order.Type)
	if err != nil {
		return nil, err
	}

	return model.CreateOrder(
		db,
		orderType,
		conf.Order.Symbol,
		conf.Order.Size,
		conf.Order.Price,
	)
}
