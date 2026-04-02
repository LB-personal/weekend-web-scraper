package worker

import (
	"bytes"

	"github.com/LB-personal/weekend-web-scraper/internal/domain"
	html "golang.org/x/net/html"
)

type Report struct {
	nextUrl     string
	currentPage int
	isLast      bool
}

func NewBasicWorker(in <-chan []byte, report <-chan Report) <-chan domain.Book {
	o := make(chan domain.Book)
	go parse(in, report, o)
	return o
}

func parse(in <-chan []byte, report <-chan Report, out chan<- domain.Book) {
	for body := range in {
		tokenizer := html.NewTokenizer(bytes.NewReader(body))
		print(tokenizer)
	}
}
