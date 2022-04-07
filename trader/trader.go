package trader

import (
	"context"

	"github.com/omhen/swissblock-trade-executor/v2/client"
	"github.com/omhen/swissblock-trade-executor/v2/model"
)

type Trader interface {
	MatchOrder(ctx context.Context, order *model.Order, reader client.BookReader) error
}
