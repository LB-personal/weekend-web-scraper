package fatcher

import (
	"context"
	"net/http"
)

type HttpGetter interface {
	Get(url string) (resp *http.Response, err error)
}

type Fatcher struct {
	client HttpGetter
	urls   <-chan string
}

func (f Fatcher) Fatch(ctx context.Context) chan<- string {
	return nil
}
