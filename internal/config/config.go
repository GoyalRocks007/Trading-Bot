package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Risk struct {
	RiskPerTradePct float64 `yaml:"risk_per_trade_pct"`
	MaxDailyLossPct float64 `yaml:"max_daily_loss_pct"`
}

type ORB struct {
	Start         string  `yaml:"start"`
	End           string  `yaml:"end"`
	RR            float64 `yaml:"rr"`
	SquareOffTime string  `yaml:"square_off_time"`
}

type Broker struct {
	Name        string `yaml:"name"`
	APIKey      string `yaml:"api_key"`
	APISecret   string `yaml:"api_secret"`
	AccessToken string `yaml:"access_token"`
}

type Logging struct {
	TradesCSV string `yaml:"trades_csv"`
	EquityCSV string `yaml:"equity_csv"`
}

type Config struct {
	Mode    string   `yaml:"mode"`
	Risk    Risk     `yaml:"risk"`
	ORB     ORB      `yaml:"orb"`
	Symbols []string `yaml:"symbols"`
	Broker  Broker   `yaml:"broker"`
	Logging Logging  `yaml:"logging"`
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
