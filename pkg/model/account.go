package model

type Settings struct {
	Pairs []string
}

type Account struct {
	Balances []Balance
}

type Balance struct {
	Asset    string
	Free     float64
	Lock     float64
	Leverage float64
}
