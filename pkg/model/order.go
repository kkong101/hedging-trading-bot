package model

import (
	"github.com/adshao/go-binance/v2/futures"
	"time"
)

type Order struct {
	ID         int64                   `json:"id"`
	ExchangeID int64                   `json:"exchange_id"`
	Pair       string                  `json:"pair"`
	Side       futures.SideType        `json:"side"`
	Type       futures.OrderType       `json:"type"`
	Status     futures.OrderStatusType `json:"status"`
	Price      float64                 `json:"price"`
	Quantity   float64                 `json:"quantity"`
	CreatedAt  time.Time               `json:"created_at"`
	UpdatedAt  time.Time               `json:"updated_at"`
	Stop       *float64                `json:"stop"`
	GroupID    *int64                  `json:"group_id"`
	RefPrice   float64                 `json:"ref_price"`
	Profit     float64                 `json:"profit"`
	Candle     Candle                  `json:"-"`
}
