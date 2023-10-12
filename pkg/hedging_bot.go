package pkg

import (
	"context"
	"modelH/pkg/egine/exchange"
	"modelH/pkg/egine/order"
	_service "modelH/pkg/egine/serivce"
	"modelH/pkg/model"
)

type CandleSubscriber interface {
	OnCandle(model.Candle)
}

type HedgingBot struct {
	exchange _service.Exchange

	settings            model.Settings
	OrderController     *order.Controller
	priorityQueueCandle *model.PriorityQueue
	OrderFeed           *order.Feed
	dataFeed            *exchange.DataFeedSubscription
}

func NewHedgingBot(settings model.Settings, e _service.Exchange) *HedgingBot {
	return &HedgingBot{
		exchange:            e,
		settings:            settings,
		OrderFeed:           order.NewOrderFeed(),
		dataFeed:            exchange.NewDataFeed(e),
		priorityQueueCandle: model.NewPriorityQueue(nil),
	}
}

func (b *HedgingBot) SubscribeCandle(subscriptions ...CandleSubscriber) {
	for _, pair := range b.settings.Pairs {
		for _, subscription := range subscriptions {
			b.dataFeed.Subscribe(pair, "1m", subscription.OnCandle, false)
		}
	}
}

func (b *HedgingBot) onCandle(candle model.Candle) {
	b.priorityQueueCandle.Push(candle)
}

func (b *HedgingBot) Run(ctx context.Context) error {

	for _, pair := range b.settings.Pairs {

		// link to ninja bot controller
		b.dataFeed.Subscribe(pair, "1m", b.onCandle, false)
	}

	// start order feed and controller
	b.OrderFeed.Start()
	b.OrderController.Start()
	defer b.OrderController.Stop()

	// start data feed and receives new candles
	b.dataFeed.Start(true)

	return nil
}
