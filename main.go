package main

import (
	"context"
	"modelH/pkg"
	"modelH/pkg/egine/exchange"
	"modelH/pkg/egine/order"
	"modelH/pkg/model"
)

func init() {

}

func main() {

	binanceFuture := exchange.NewBinanceFuture(
		context.Background(),
		exchange.WithBinanceFutureHedgeMode(true),
		exchange.WithBinanceFutureCredentials("", ""),
	)

	settings := model.Settings{Pairs: []string{"BTCUSDT", "ETHUSDT"}}
	bot := pkg.NewHedgingBot(settings, binanceFuture)

	bot.OrderController = order.NewController(context.Background(), binanceFuture, bot.OrderFeed)

	err := bot.Run(context.Background())
	if err != nil {
		panic(err)
	}

}
