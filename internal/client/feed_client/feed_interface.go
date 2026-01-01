package feedclient

import (
	"context"
	"net/http"
	"sync"
)

type IFeedClient interface {
	StartAuth(ctx context.Context) (string, error)
	HandleCallback(ctx context.Context, r *http.Request) error
	Start(wg *sync.WaitGroup) error
	Stop() error
}
