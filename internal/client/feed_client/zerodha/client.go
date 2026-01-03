package zerodhafeedclient

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
	"trading-bot/internal/models"
	"trading-bot/logger"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
	kitemodels "github.com/zerodha/gokiteconnect/v4/models"
	kiteticker "github.com/zerodha/gokiteconnect/v4/ticker"
	"go.uber.org/zap"
)

var (
	api_key      = ""
	api_secret   = ""
	access_token = ""

	zerodhaFeedClient *ZerodhaFeedClient

	GetZerodhaFeedClient = func(bus *models.Bus) *ZerodhaFeedClient {
		api_key = os.Getenv("KITE_API_KEY")
		api_secret = os.Getenv("KITE_API_SECRET")
		if zerodhaFeedClient == nil {

			kc := kiteconnect.New(api_key)

			zerodhaFeedClient = &ZerodhaFeedClient{
				kc:  kc,
				bus: bus,
			}
		}
		return zerodhaFeedClient
	}
)

type ZerodhaFeedClient struct {
	kc     *kiteconnect.Client
	ticker *kiteticker.Ticker
	bus    *models.Bus
	tokens []uint32
}

func (zc *ZerodhaFeedClient) StartAuth(ctx context.Context) (string, error) {
	url := zc.kc.GetLoginURL()
	return url, nil
}

func (zc *ZerodhaFeedClient) HandleCallback(ctx context.Context, r *http.Request) error {
	requestToken := r.URL.Query().Get("request_token")
	data, err := zc.kc.GenerateSession(requestToken, api_secret)
	if err != nil {
		return err
	}
	access_token = data.AccessToken

	return nil
}

func (zc *ZerodhaFeedClient) Subscribe() error {
	tokens := []uint32{16401, 738561, 60417, 1304833, 2707457}
	zc.tokens = tokens
	if zc.ticker != nil {
		if err := zc.ticker.Subscribe(tokens); err != nil {
			return err
		}
		if err := zc.ticker.SetMode(kiteticker.ModeFull, tokens); err != nil {
			return err
		}
	}
	return nil
}

func (zc *ZerodhaFeedClient) registerCallbacks() {

	zc.ticker.OnError(func(err error) {
		logger.Log.Error("Ticker error:", zap.Error(err))
	})

	zc.ticker.OnClose(func(code int, reason string) {
		logger.Log.Info("Ticker closed:", zap.Int("code", code), zap.String("reason", reason))
	})

	zc.ticker.OnConnect(func() {
		logger.Log.Info("Ticker connected")
		// re-subscribe on connect if tokens exist
		zc.Subscribe()
	})

	zc.ticker.OnTick(func(tick kitemodels.Tick) {
		internalTick := ZerodhaTickToInternalTickMapper(tick)
		// push into channel (non-blocking best-effort)
		select {
		case zc.bus.Ticks <- internalTick:
		default:
			// channel full, drop tick or handle backpressure here
			logger.Log.Warn("tick channel full, dropping tick")
		}
		select {
		case zc.bus.PositionTicks <- internalTick:
		default:
			logger.Log.Warn("position ticks channel full, dropping tick")
		}
	})

	// Optional: handle reconnection notifications
	zc.ticker.OnReconnect(func(attempt int, delay time.Duration) {
		logger.Log.Info("Reconnect attempt", zap.Int("attempt", attempt), zap.Duration("delay", delay))
	})
	zc.ticker.OnNoReconnect(func(attempt int) {
		logger.Log.Info("No more reconnects", zap.Int("attempt", attempt))
	})

	zc.ticker.OnOrderUpdate(func(order kiteconnect.Order) {
		logger.Log.Info("Order update:", zap.String("order_id", order.OrderID))
	})
}

func (zc *ZerodhaFeedClient) Stop() error {
	if zc.ticker != nil {
		return zc.ticker.Close()
	}
	return fmt.Errorf("ticker not initialized")
}

func (zc *ZerodhaFeedClient) Init() {
	kt := kiteticker.New(api_key, access_token)
	zc.ticker = kt
}

func (zc *ZerodhaFeedClient) Start(wg *sync.WaitGroup) error {
	zc.Init()
	if zc.ticker == nil {
		return fmt.Errorf("ticker not initialized")
	}
	zc.registerCallbacks()
	wg.Add(1)
	go func() {
		defer wg.Done()
		zc.ticker.Serve()
	}()
	return nil
}
