package client

import (
	"context"

	"github.com/omhen/swissblock-trade-executor/v2/model"
)

type BookReader interface {
	StartReading(ctx context.Context, symbol string) (chan *model.BookItem, chan error)
	Stop(ctx context.Context)
}
