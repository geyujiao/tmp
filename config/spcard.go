package config

var Spcard spcard

type spcard struct {
	ChargeNumber string `config:"chargeNumber"`
	ListenPort   string `config:"listenPort"`
	ChargeType   string `config:"type"`
}
