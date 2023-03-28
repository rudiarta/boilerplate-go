package libelastic

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"syscall"
	"time"

	es "github.com/olivere/elastic/v7"
)

// ElasticRetrier ....
type ElasticRetrier struct {
	backoff     es.Backoff
	onRetryFunc func(err error)
}

// NewElasticRetrier ...
func NewElasticRetrier(t time.Duration, f func(err error)) *ElasticRetrier {
	return &ElasticRetrier{
		backoff:     es.NewConstantBackoff(t),
		onRetryFunc: f,
	}
}

// Retry ...
func (r *ElasticRetrier) Retry(ctx context.Context, retry int, req *http.Request, resp *http.Response, err error) (time.Duration, bool, error) {
	log.Print(fmt.Errorf(fmt.Sprintf("elasticsearch Retrier #%d", retry)))
	if err == syscall.ECONNREFUSED {
		err = fmt.Errorf("elasticsearch or network down")
	}
	wait, stop := r.backoff.Next(retry)
	r.onRetryFunc(err)
	return wait, stop, nil
}
