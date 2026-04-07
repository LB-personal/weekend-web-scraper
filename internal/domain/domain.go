package domain

type Book struct {
	Name     string
	Rating   Rating
	Price    Money
	Category string
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
