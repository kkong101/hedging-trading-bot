package order

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	log "github.com/sirupsen/logrus"
	"math"
	"modelH/pkg/egine/exchange"
	_service "modelH/pkg/egine/serivce"
	"modelH/pkg/model"
	"strings"
	"sync"
	"time"
)

type summary struct {
	Pair      string
	WinLong   []float64
	WinShort  []float64
	LoseLong  []float64
	LoseShort []float64
	Volume    float64
}

func (s summary) Win() []float64 {
	return append(s.WinLong, s.WinShort...)
}

func (s summary) Lose() []float64 {
	return append(s.LoseLong, s.LoseShort...)
}

func (s summary) Profit() float64 {
	profit := 0.0
	for _, value := range append(s.Win(), s.Lose()...) {
		profit += value
	}
	return profit
}

func (s summary) SQN() float64 {
	total := float64(len(s.Win()) + len(s.Lose()))
	avgProfit := s.Profit() / total
	stdDev := 0.0
	for _, profit := range append(s.Win(), s.Lose()...) {
		stdDev += math.Pow(profit-avgProfit, 2)
	}
	stdDev = math.Sqrt(stdDev / total)
	return math.Sqrt(total) * (s.Profit() / total) / stdDev
}

func (s summary) Payoff() float64 {
	avgWin := 0.0
	avgLose := 0.0

	for _, value := range s.Win() {
		avgWin += value
	}

	for _, value := range s.Lose() {
		avgLose += value
	}

	if len(s.Win()) == 0 || len(s.Lose()) == 0 || avgLose == 0 {
		return 0
	}

	return (avgWin / float64(len(s.Win()))) / math.Abs(avgLose/float64(len(s.Lose())))
}

func (s summary) WinPercentage() float64 {
	if len(s.Win())+len(s.Lose()) == 0 {
		return 0
	}

	return float64(len(s.Win())) / float64(len(s.Win())+len(s.Lose())) * 100
}

func (s summary) String() string {
	tableString := &strings.Builder{}

	return tableString.String()
}

type Status string

const (
	StatusRunning Status = "running"
	StatusStopped Status = "stopped"
	StatusError   Status = "error"
)

type Controller struct {
	mtx            sync.Mutex
	ctx            context.Context
	exchange       _service.Exchange
	orderFeed      *Feed
	Results        map[string]*summary
	lastPrice      map[string]float64
	tickerInterval time.Duration
	finish         chan bool
	status         Status
}

func NewController(ctx context.Context, exchange _service.Exchange,
	orderFeed *Feed) *Controller {

	return &Controller{
		ctx:            ctx,
		exchange:       exchange,
		orderFeed:      orderFeed,
		lastPrice:      make(map[string]float64),
		Results:        make(map[string]*summary),
		tickerInterval: time.Second,
		finish:         make(chan bool),
	}
}

func (c *Controller) OnCandle(candle model.Candle) {
	c.lastPrice[candle.Pair] = candle.Close
}

func (c *Controller) calculateProfit(o *model.Order) (value, percent float64, err error) {
	// get filled orders before the current order

	quantity := 0.0
	avgPriceLong := 0.0
	avgPriceShort := 0.0

	fmt.Println(quantity, avgPriceLong, avgPriceShort)

	return 0, 0, nil
}

func (c *Controller) processTrade(order *model.Order) {
	if order.Status != futures.OrderStatusTypeFilled {
		return
	}

	// initializer results map if needed
	if _, ok := c.Results[order.Pair]; !ok {
		c.Results[order.Pair] = &summary{Pair: order.Pair}
	}

	// register order volume
	c.Results[order.Pair].Volume += order.Price * order.Quantity

	profitValue, profit, err := c.calculateProfit(order)
	if err != nil {
		return
	}

	order.Profit = profit
	if profitValue == 0 {
		return
	} else if profitValue > 0 {
		if order.Side == futures.SideTypeBuy {
			c.Results[order.Pair].WinLong = append(c.Results[order.Pair].WinLong, profitValue)
		} else {
			c.Results[order.Pair].WinShort = append(c.Results[order.Pair].WinShort, profitValue)
		}
	} else {
		if order.Side == futures.SideTypeBuy {
			c.Results[order.Pair].LoseLong = append(c.Results[order.Pair].LoseLong, profitValue)
		} else {
			c.Results[order.Pair].LoseShort = append(c.Results[order.Pair].LoseShort, profitValue)
		}
	}

	_, _ = exchange.SplitAssetQuote(order.Pair)
}

func (c *Controller) updateOrders() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

}

func (c *Controller) Status() Status {
	return c.status
}

func (c *Controller) Start() {
	if c.status != StatusRunning {
		c.status = StatusRunning
		go func() {
			ticker := time.NewTicker(c.tickerInterval)
			for {
				select {
				case <-ticker.C:
					c.updateOrders()
				case <-c.finish:
					ticker.Stop()
					return
				}
			}
		}()
		log.Info("Bot started.")
	}
}

func (c *Controller) Stop() {
	if c.status == StatusRunning {
		c.status = StatusStopped
		c.updateOrders()
		c.finish <- true
		log.Info("Bot stopped.")
	}
}

func (c *Controller) GetAccount() (model.Account, error) {
	return c.exchange.GetAccount()
}

func (c *Controller) GetPosition(pair string) (asset, quote float64, err error) {
	return c.exchange.GetPosition(pair)
}

func (c *Controller) LastQuote(pair string) (float64, error) {
	return c.exchange.LastQuote(c.ctx, pair)
}

func (c *Controller) PositionValue(pair string) (float64, error) {
	asset, _, err := c.exchange.GetPosition(pair)
	if err != nil {
		return 0, err
	}
	return asset * c.lastPrice[pair], nil
}

func (c *Controller) GetOrder(pair string, id int64) (model.Order, error) {
	return c.exchange.GetOrder(pair, id)
}

func (c *Controller) CreateOrderOCO(side futures.SideType, pair string, size, price, stop,
	stopLimit float64) ([]model.Order, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return nil, nil
}

func (c *Controller) OpenLimitOrder(side futures.SideType, pair string, size, limit float64) (model.Order, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	log.Infof("[ORDER] Creating LIMIT %s order for %s", side, pair)
	order, err := c.exchange.OpenLimitOrder(side, pair, size, limit)
	if err != nil {
		return model.Order{}, err
	}

	go c.orderFeed.Publish(order, true)
	log.Infof("[ORDER CREATED] %s", order)
	return order, nil
}

func (c *Controller) CreateOrderMarketQuote(side futures.SideType, pair string, amount float64) (model.Order, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return model.Order{}, nil
}

func (c *Controller) OpenMarketOrder(side futures.SideType, pair string, size float64) (model.Order, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	log.Infof("[ORDER] Creating MARKET %s order for %s", side, pair)
	order, err := c.exchange.OpenMarketOrder(side, pair, size)
	if err != nil {
		return model.Order{}, err
	}

	// calculate profit
	c.processTrade(&order)
	go c.orderFeed.Publish(order, true)
	log.Infof("[ORDER CREATED] %s", order)
	return order, err
}

func (c *Controller) CancelOpenOrder(pair string, size float64, limit float64) (model.Order, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	log.Infof("[ORDER] Creating STOP order for %s", pair)
	order, err := c.exchange.CancelOpenOrder(pair, size, limit)
	if err != nil {
		return model.Order{}, err
	}

	go c.orderFeed.Publish(order, true)
	log.Infof("[ORDER CREATED] %s", order)
	return order, nil
}

func (c *Controller) Cancel(order model.Order) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return nil
}
