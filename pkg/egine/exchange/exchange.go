package exchange

import (
	"context"
	"fmt"
	"github.com/StudioSol/set"
	"log"
	_service "modelH/pkg/egine/serivce"
	"modelH/pkg/model"
	"strings"
	"sync"
)

type DataFeedSubscription struct {
	exchange                _service.Exchange
	Feeds                   *set.LinkedHashSetString
	DataFeeds               map[string]*DataFeed
	SubscriptionsByDataFeed map[string][]Subscription
}

type DataFeed struct {
	Data chan model.Candle
	Err  chan error
}

type Subscription struct {
	onCandleClose bool
	consumer      DataFeedConsumer
}

type DataFeedConsumer func(model.Candle)

func NewDataFeed(exchange _service.Exchange) *DataFeedSubscription {
	return &DataFeedSubscription{
		exchange:                exchange,
		Feeds:                   set.NewLinkedHashSetString(),
		DataFeeds:               make(map[string]*DataFeed),
		SubscriptionsByDataFeed: make(map[string][]Subscription),
	}
}

func (d *DataFeedSubscription) Connect() {
	log.Println("Connecting to the exchange.")
	for feed := range d.Feeds.Iter() {
		pair, timeframe := d.pairTimeframeFromKey(feed)
		candleChan, errChan := d.exchange.CandlesSubscription(context.Background(), pair, timeframe)
		d.DataFeeds[feed] = &DataFeed{
			Data: candleChan,
			Err:  errChan,
		}
	}
}

func (d *DataFeedSubscription) Start(loadSync bool) {
	d.Connect()
	wg := new(sync.WaitGroup)

	for key, feed := range d.DataFeeds {
		wg.Add(1)
		go func(key string, feed *DataFeed) {
			for {
				select {
				case candle, ok := <-feed.Data:
					if !ok {
						wg.Done()
						return
					}
					for _, subscription := range d.SubscriptionsByDataFeed[key] {
						if subscription.onCandleClose && !candle.Complete {
							continue
						}
						subscription.consumer(candle)
					}
				case err := <-feed.Err:
					if err != nil {
						log.Println(err)
					}
				}
			}
		}(key, feed)
	}

	log.Println("Data feed connected.")
	if loadSync {
		wg.Wait()
	}
}

func (d *DataFeedSubscription) Subscribe(pair, timeframe string, consumer DataFeedConsumer, onCandleClose bool) {
	key := d.feedKey(pair, timeframe)
	d.Feeds.Add(key)

	d.SubscriptionsByDataFeed[key] = append(d.SubscriptionsByDataFeed[key], Subscription{
		onCandleClose: onCandleClose,
		consumer:      consumer,
	})
}

func (d *DataFeedSubscription) feedKey(pair, timeframe string) string {
	return fmt.Sprintf("%s--%s", pair, timeframe)
}

func (d *DataFeedSubscription) pairTimeframeFromKey(key string) (pair, timeframe string) {
	parts := strings.Split(key, "--")
	return parts[0], parts[1]
}
