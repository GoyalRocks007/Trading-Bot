package zerodhafeedclient

import (
	"trading-bot/internal/models"

	kitemodels "github.com/zerodha/gokiteconnect/v4/models"
)

func ZerodhaTickToInternalTickMapper(tick kitemodels.Tick) models.Tick {
	return models.Tick{
		Symbol:             getSymbolFromZerodhaToken(tick.InstrumentToken),
		LTP:                tick.LastPrice,
		LastTradedQuantity: int32(tick.LastTradedQuantity),
		VolumeTraded:       int32(tick.VolumeTraded),
		Time:               tick.Timestamp.Time,
	}
}

func getSymbolFromZerodhaToken(token uint32) string {
	switch token {
	case 6401:
		return "ADANIENT"
	case 738561:
		return "RELIANCE"
	case 2707457:
		return "NIFTYBEES"
	case 60417:
		return "ASIANPAINT"
	case 1304833:
		return "ETERNAL"
	default:
		return ""
	}
}
