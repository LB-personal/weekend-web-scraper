package domain

import "net/url"

type Book struct {
	Name   string
	Rating Rating
	Price  Money
	//	Category   string - requires enrichment, will skip for not and add after mvp
	InStock    bool
	DetailPage *url.URL
}

type Stats struct {
	TotalBooks    int
	Categories    int
	MostExpensive Book
	Cheapest      Book
	AvaragePrice  Money
	CountByRating map[Rating]int
}

type PageData struct {
	Content []byte
	Url     string
}
