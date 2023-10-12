package serivce

import (
	"context"
	"modelH/pkg/model"
	"time"
)

type Exchange interface {
	Broker
	Feeder
}

// Feeder is the interface that provides data for the strategy.
type Feeder interface {
	GetPairInfo(pair string) model.Pair
	LastQuote(ctx context.Context, pair string) (float64, error)
	CandlesByPeriod(ctx context.Context, pair, period string, start, end time.Time) ([]model.Candle, error)
	CandlesByLimit(ctx context.Context, pair, period string, limit int) ([]model.Candle, error)
	CandlesSubscription(ctx context.Context, pair, timeframe string) (chan model.Candle, chan error)
}

type Broker interface {
	GetAccount() (model.Account, error)
	GetPosition(pair string) (asset, quote float64, err error)
	GetAllPositions() ([]model.Position, error)

	GetOrder(pair string, id int64) (model.Order, error)

	OpenLimitOrder(side model.SideType, pair string, size float64, limit float64) (model.Order, error)
	OpenMarketOrder(side model.SideType, pair string, size float64) (model.Order, error)

	CloseBuyPosition(pair string, size float64) (model.Order, error)
	CloseSellPosition(pair string, size float64) (model.Order, error)

	CancelOpenOrder(pair string, quantity float64, limit float64) (model.Order, error)
	CancelAllOpenOrder(pair string) ([]model.Order, error)
}

type Notifier interface {
	Notify(string)
	OnOrder(order model.Order)
	OnError(err error)
}
