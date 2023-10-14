package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/jpillora/backoff"
	"log"
	"modelH/pkg/model"
	"strconv"
	"time"
)

type PairOption struct {
	Pair       string
	Leverage   int
	MarginType futures.MarginType
}

type BinanceFuture struct {
	ctx context.Context

	APIKey    string
	APISecret string

	client      *futures.Client
	pairs       map[string]model.Pair
	isTestNet   bool
	isHedgeMode bool
}

func (b *BinanceFuture) GetAllPositions() ([]model.Position, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BinanceFuture) CloseBuyPosition(pair string, size float64) (model.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BinanceFuture) CloseSellPosition(pair string, size float64) (model.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BinanceFuture) CancelAllOpenOrder(pair string) ([]model.Order, error) {
	//TODO implement me
	panic("implement me")
}

type BinanceFutureOption func(*BinanceFuture)

func NewBinanceFuture(ctx context.Context, options ...BinanceFutureOption) *BinanceFuture {

	binance.WebsocketKeepalive = true
	exchange := &BinanceFuture{ctx: ctx}

	for _, option := range options {
		option(exchange)
	}

	exchange.client = futures.NewClient(exchange.APIKey, exchange.APISecret)
	err := exchange.client.NewPingService().Do(ctx)
	if err != nil {
		panic(err)
	}

	results, err := exchange.client.NewExchangeInfoService().Do(ctx)
	if err != nil {
		panic(err)
	}

	// 모든 포지션을 정리해야지 양방향 모드 가능함.
	err = exchange.client.NewChangePositionModeService().DualSide(exchange.isHedgeMode).Do(ctx)
	if err != nil {
		panic(err)
	}

	// Initialize with orders precision and assets limits
	exchange.pairs = make(map[string]model.Pair)
	for _, info := range results.Symbols {
		tradeLimits := model.Pair{
			BaseAsset:          info.BaseAsset,
			QuoteAsset:         info.QuoteAsset,
			BaseAssetPrecision: info.BaseAssetPrecision,
			QuotePrecision:     info.QuotePrecision,
		}
		for _, filter := range info.Filters {
			if typ, ok := filter["filterType"]; ok {
				if typ == string(binance.SymbolFilterTypeLotSize) {
					tradeLimits.MinQuantity, _ = strconv.ParseFloat(filter["minQty"].(string), 64)
					tradeLimits.MaxQuantity, _ = strconv.ParseFloat(filter["maxQty"].(string), 64)
					tradeLimits.StepSize, _ = strconv.ParseFloat(filter["stepSize"].(string), 64)
				}

				if typ == string(binance.SymbolFilterTypePriceFilter) {
					tradeLimits.MinPrice, _ = strconv.ParseFloat(filter["minPrice"].(string), 64)
					tradeLimits.MaxPrice, _ = strconv.ParseFloat(filter["maxPrice"].(string), 64)
					tradeLimits.TickSize, _ = strconv.ParseFloat(filter["tickSize"].(string), 64)
				}
			}
		}
		exchange.pairs[info.Symbol] = tradeLimits
	}

	for _, v := range exchange.pairs {
		re, _ := json.MarshalIndent(v, "", "  ")
		log.Println(string(re))
	}

	log.Println("[SETUP] Using Binance Futures exchange")

	return exchange
}

func (b *BinanceFuture) GetAccount() (model.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BinanceFuture) GetPosition(pair string) (asset, quote float64, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *BinanceFuture) GetOrder(pair string, id int64) (model.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BinanceFuture) CreateOrderOCO(side futures.SideType, pair string, size, price, stop, stopLimit float64) ([]model.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BinanceFuture) OpenLimitOrder(side futures.SideType, pair string, size float64, limit float64) (model.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BinanceFuture) OpenMarketOrder(side futures.SideType, pair string, size float64) (model.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BinanceFuture) CreateOrderMarketQuote(side futures.SideType, pair string, quote float64) (model.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BinanceFuture) CancelOpenOrder(pair string, quantity float64, limit float64) (model.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BinanceFuture) Cancel(order model.Order) error {
	//TODO implement me
	panic("implement me")
}

func (b *BinanceFuture) PairInfo(pair string) model.Pair {
	//TODO implement me
	panic("implement me")
}

func (b *BinanceFuture) LastQuote(ctx context.Context, pair string) (float64, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BinanceFuture) CandlesByPeriod(ctx context.Context, pair, period string, start, end time.Time) ([]model.Candle, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BinanceFuture) CandlesByLimit(ctx context.Context, pair, period string, limit int) ([]model.Candle, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BinanceFuture) CandlesSubscription(ctx context.Context, pair, period string) (chan model.Candle, chan error) {
	candleChan := make(chan model.Candle)
	errChan := make(chan error)

	go func() {
		ba := &backoff.Backoff{
			Min: 100 * time.Millisecond,
			Max: 1 * time.Second,
		}

		for {
			done, _, err := futures.WsKlineServe(pair, period, func(event *futures.WsKlineEvent) {
				ba.Reset()

				candle := toCandleFromWsKline(pair, event.Kline)

				// !@#!@#
				tt, _ := json.MarshalIndent(candle, "", "  ")
				fmt.Println(string(tt))

				if candle.Complete {
					log.Println("Complete candle")
				}

				candleChan <- candle

			}, func(err error) {
				errChan <- err
			})
			if err != nil {
				errChan <- err
				close(errChan)
				close(candleChan)
				return
			}

			select {
			case <-ctx.Done():
				close(errChan)
				close(candleChan)
				return
			case <-done:
				time.Sleep(ba.Duration())
			}
		}
	}()

	return candleChan, errChan
}

func WithBinanceFutureTestNet() BinanceFutureOption {
	return func(b *BinanceFuture) {
		b.isTestNet = true
	}
}

func WithBinanceFutureHedgeMode(isHedgeMode bool) BinanceFutureOption {
	return func(b *BinanceFuture) {
		b.isHedgeMode = isHedgeMode
	}
}

func WithBinanceFutureCredentials(key, secret string) BinanceFutureOption {
	return func(b *BinanceFuture) {
		b.APIKey = key
		b.APISecret = secret
	}
}

func (b *BinanceFuture) Run() {

}

func toCandleFromKline(pair string, k futures.Kline) model.Candle {
	var err error
	t := time.Unix(0, k.OpenTime*int64(time.Millisecond))
	candle := model.Candle{Pair: pair, Time: t, UpdatedAt: t}
	candle.Open, err = strconv.ParseFloat(k.Open, 64)
	candle.Close, err = strconv.ParseFloat(k.Close, 64)
	candle.High, err = strconv.ParseFloat(k.High, 64)
	candle.Low, err = strconv.ParseFloat(k.Low, 64)
	candle.Volume, err = strconv.ParseFloat(k.Volume, 64)
	candle.Complete = true
	candle.Metadata = make(map[string]float64)
	if err != nil {
		log.Println(err)
	}
	return candle
}

func toCandleFromWsKline(pair string, k futures.WsKline) model.Candle {
	var err error
	t := time.Unix(0, k.StartTime*int64(time.Millisecond))
	candle := model.Candle{Pair: pair, Time: t, UpdatedAt: t}
	candle.Open, err = strconv.ParseFloat(k.Open, 64)
	candle.Close, err = strconv.ParseFloat(k.Close, 64)
	candle.High, err = strconv.ParseFloat(k.High, 64)
	candle.Low, err = strconv.ParseFloat(k.Low, 64)
	candle.Volume, err = strconv.ParseFloat(k.Volume, 64)
	candle.Complete = k.IsFinal
	candle.Metadata = make(map[string]float64)
	if err != nil {
		log.Println(err)
	}
	return candle
}
