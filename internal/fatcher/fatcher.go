package fatcher

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
)

type HttpGetter interface {
	Get(url string) (resp *http.Response, err error)
}

type fatcher struct {
	client HttpGetter
	urls   <-chan string
}

type Fatcher interface {
	Fatch(ctx context.Context) <-chan []byte
}

func NewFatcher(c HttpGetter, u <-chan string) Fatcher {
	return &fatcher{
		client: c,
		urls:   u,
	}
}

func (f *fatcher) Fatch(ctx context.Context) <-chan []byte {
	c := make(chan []byte)
	go f.internal(ctx, c)
	return c
}

func (f *fatcher) internal(ctx context.Context, responses chan<- []byte) {
	for {
		select {
		case <-ctx.Done():
			return
		case url, ok := <-f.urls:
			if !ok {
				return
			}

			res, err := f.client.Get(url)
			if err != nil {
				log.Fatal(err)
			}

			body, err := handleRespone(res)
			if err != nil {
				log.Fatal(err)
			}

			responses <- body
		}
	}
}

func handleRespone(res *http.Response) ([]byte, error) {
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, errors.New("got bad request")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
