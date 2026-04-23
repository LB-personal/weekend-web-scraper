package worker_test

import (
	_ "embed"
	"net/url"
	"testing"

	"github.com/LB-personal/weekend-web-scraper/internal/domain"
	"github.com/LB-personal/weekend-web-scraper/internal/worker"
)

//go:embed sample_page.html
var b []byte

func Test_parser_complete(t *testing.T) {
	in := make(chan domain.PageData)
	r := make(chan string)
	out := worker.NewParser(in, r)

	in <- domain.PageData{
		Content: b,
		Url:     "https://www.mock-base.com",
	}

	close(in)
	u := <-r
	e := "https://www.mock-base.com/catalogue/page-4.html"
	if u != e {
		t.Errorf("url: expected %s, actual %s", e, u)
	}

	book := <-out

	eb := domain.Book{
		Name:       "Slow States of Collapse: Poems",
		Rating:     must(domain.NewRateing(3)),
		Price:      must(domain.NewMoney("£57.31")),
		InStock:    true,
		DetailPage: must(url.Parse("https://www.mock-base.com/catalogue/slow-states-of-collapse-poems_960/index.html")),
	}

	if !book.Equals(eb) {
		t.Errorf("expected %v, actual %v", eb, book)
	}

	_, ok := <-out

	if ok {
		t.Error("expect channel to be closed, channel was open")
	}
}

func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}

	return val
}
