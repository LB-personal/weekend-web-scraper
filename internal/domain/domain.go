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

func (b Book) Equals(o Book) bool {
	return b.Name == o.Name &&
		b.Rating == o.Rating &&
		b.Price == o.Price &&
		b.InStock == o.InStock &&
		*b.DetailPage == *o.DetailPage
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
