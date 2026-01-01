package models

import "time"

type Tick struct {
	Symbol             string
	LTP                float64
	LastTradedQuantity int32
	VolumeTraded       int32
	Time               time.Time
}

type Candle struct {
	Symbol   string
	Start    time.Time
	End      time.Time
	Duration int
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   int64
}
