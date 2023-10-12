package model

import (
	"github.com/adshao/go-binance/v2/futures"
)

/**

// PositionRisk define position risk info
	type PositionRisk struct {
		EntryPrice       string `json:"entryPrice"`
		MarginType       string `json:"marginType"`
		IsAutoAddMargin  string `json:"isAutoAddMargin"`
		IsolatedMargin   string `json:"isolatedMargin"`
		Leverage         string `json:"leverage"`
		LiquidationPrice string `json:"liquidationPrice"`
		MarkPrice        string `json:"markPrice"`
		MaxNotionalValue string `json:"maxNotionalValue"`
		PositionAmt      string `json:"positionAmt"`
		Symbol           string `json:"symbol"`
		UnRealizedProfit string `json:"unRealizedProfit"`
		PositionSide     string `json:"positionSide"`
		Notional         string `json:"notional"`
		IsolatedWallet   string `json:"isolatedWallet"`
	}
*/

type Position struct {
	Pair        string
	Side        futures.SideType
	PositionAmt float64
	EntryPrice  float64
	Leverage    int
	MarginType  futures.MarginType
}
