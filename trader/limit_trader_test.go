package trader_test

import (
	"context"
	"time"

	"github.com/omhen/swissblock-trade-executor/v2/client"
	"github.com/omhen/swissblock-trade-executor/v2/database"
	"github.com/omhen/swissblock-trade-executor/v2/errors"
	"github.com/omhen/swissblock-trade-executor/v2/model"
	"github.com/omhen/swissblock-trade-executor/v2/trader"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockReader struct {
	messages []*model.BookItem
	errors   []error
}

func NewMockReader(messages []*model.BookItem, errors []error) client.BookReader {
	return &mockReader{
		messages: messages,
		errors:   errors,
	}
}

func (mr *mockReader) StartReading(ctx context.Context, symbol string) (chan *model.BookItem, chan error) {
	receiver := make(chan *model.BookItem)
	errChannel := make(chan error)

	go func() {
		for _, err := range mr.errors {
			errChannel <- err
		}

		for _, message := range mr.messages {
			receiver <- message
		}
	}()

	return receiver, errChannel
}

func (mr *mockReader) Stop(ctx context.Context) {
}

var _ = Describe("LimitTrader", func() {
	db, _ := database.NewDBConnection("sqlite:///tmp/test")
	ctx := context.WithValue(context.Background(), database.ContextKeyDB, db)
	ordersInBook := []*model.BookItem{
		{UpdateID: 1, Symbol: "BNBUSDT", BestBidPrice: 90, BestBidQuantity: 9, BestAskPrice: 95, BestAskQuantity: 9},
		{UpdateID: 1, Symbol: "BNBUSDT", BestBidPrice: 80, BestBidQuantity: 10, BestAskPrice: 90, BestAskQuantity: 10},
		{UpdateID: 1, Symbol: "BNBUSDT", BestBidPrice: 100, BestBidQuantity: 10, BestAskPrice: 110, BestAskQuantity: 10},
	}

	Describe("Executing a sell order", func() {
		Context("and there are no errors", func() {
			It("should generate a single match", func() {
				reader := NewMockReader(ordersInBook, []error{})
				order, err := model.CreateOrder(db, model.Sell, "BNBUSDT", 10, 100)
				Expect(err).NotTo(HaveOccurred())
				trader := trader.NewLimitTrader(1 * time.Second)
				err = trader.MatchOrder(ctx, order, reader)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(order.Trades)).To(Equal(1))
			})
			It("should generate a multiple match", func() {
				reader := NewMockReader(ordersInBook, []error{})
				order, err := model.CreateOrder(db, model.Sell, "BNBUSDT", 10, 90)
				Expect(err).NotTo(HaveOccurred())
				trader := trader.NewLimitTrader(1 * time.Second)
				err = trader.MatchOrder(ctx, order, reader)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(order.Trades)).To(Equal(2))
			})
			It("should timeout with no matches", func() {
				reader := NewMockReader(ordersInBook, []error{})
				order, err := model.CreateOrder(db, model.Sell, "BNBUSDT", 10, 110)
				Expect(err).NotTo(HaveOccurred())
				trader := trader.NewLimitTrader(0 * time.Second)
				err = trader.MatchOrder(ctx, order, reader)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(order.Trades)).To(Equal(0))
			})
		})
	})

	Describe("Executing a buy order", func() {
		Context("and there are no errors", func() {
			It("should generate a single match", func() {
				reader := NewMockReader(ordersInBook, []error{})
				order, err := model.CreateOrder(db, model.Buy, "BNBUSDT", 10, 94)
				Expect(err).NotTo(HaveOccurred())
				trader := trader.NewLimitTrader(1 * time.Second)
				err = trader.MatchOrder(ctx, order, reader)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(order.Trades)).To(Equal(1))
			})
			It("should generate a multiple match", func() {
				reader := NewMockReader(ordersInBook, []error{})
				order, err := model.CreateOrder(db, model.Buy, "BNBUSDT", 10, 100)
				Expect(err).NotTo(HaveOccurred())
				trader := trader.NewLimitTrader(1 * time.Second)
				err = trader.MatchOrder(ctx, order, reader)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(order.Trades)).To(Equal(2))
			})
			It("should timeout with no matches", func() {
				reader := NewMockReader(ordersInBook, []error{})
				order, err := model.CreateOrder(db, model.Buy, "BNBUSDT", 10, 85)
				Expect(err).NotTo(HaveOccurred())
				trader := trader.NewLimitTrader(0 * time.Second)
				err = trader.MatchOrder(ctx, order, reader)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(order.Trades)).To(Equal(0))
			})
		})
	})

	Describe("Executing any order", func() {
		Context("and an error is returned by the reader", func() {
			It("should return the error", func() {
				reader := NewMockReader(ordersInBook, []error{errors.New("An error produced at the reader")})
				order, err := model.CreateOrder(db, model.Sell, "BNBUSDT", 10, 150)
				Expect(err).NotTo(HaveOccurred())
				trader := trader.NewLimitTrader(1 * time.Second)
				err = trader.MatchOrder(ctx, order, reader)
				Expect(err.Error()).To(Equal("While processing order: An error produced at the reader"))
			})

		})
	})
})
