package ingressmodule

import (
	"sync"
	"time"
	"trading-bot/internal/models"
)

type IIngressModule interface {
	Runner(bus *models.Bus, interval time.Duration, wg *sync.WaitGroup) error
}
