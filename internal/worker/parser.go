package worker

import (
	"log"
	"net/url"
	"strings"

	"github.com/LB-personal/weekend-web-scraper/internal/domain"
	html "golang.org/x/net/html"
)

func NewParser(in <-chan domain.PageData, report chan<- string) <-chan domain.Book {
	o := make(chan domain.Book)
	go parse(in, report, o)
	return o
}

func parse(in <-chan domain.PageData, report chan<- string, out chan<- domain.Book) {
	for pd := range in {
		u, l := getNextPageUrl(pd.Content)
		if l {
			close(report)
		} else {
			u := resolveUrl(pd.Url, u)
			report <- u.String()
		}

		parseBooks(pd, out)
	}

	close(out)
}

func resolveUrl(bUrl string, hrefUrl string) *url.URL {
	base, _ := url.Parse(bUrl) // cannot err, if the url is malformed the fetcher could not have fetch the content
	href, err := url.Parse(hrefUrl)

	if err != nil {
		log.Fatal("scraped site is broken", err)
	}

	u := base.ResolveReference(href)
	return u
}

func getNextPageUrl(b []byte) (next string, isLast bool) {
	s := newTokenStream(b)
	_, err := s.skipToBeginning("ul", map[string]string{"class": "pager"})
	if err != nil {
		log.Fatalf("page is missing a pager")
	}

	if _, err := s.skipToBeginning("li", map[string]string{"class": "next"}); err == nil {
		attr, err := s.skipToBeginning("a", map[string]string{})
		if err != nil {
			log.Fatal("next button is missing a next page link")
		}

		if href, ok := attr["href"]; ok {
			return href, false
		}
	}

	return "", true
}

func parseBooks(pd domain.PageData, o chan<- domain.Book) {
	s := newTokenStream(pd.Content)
	_, err := s.skipToBeginning("ol", map[string]string{"class": "row"})
	if err != nil {
		log.Fatal("size is malformed, no books to scrape")
	}

	for {
		switch s.Next() {
		case html.TextToken:
			if strings.Trim(string(s.z.Text()), "\t\r\n ") != "" {
				log.Fatalf("html is malformed")
			}
		case html.EndTagToken:
			if name, _ := s.Tag(); name != "ol" {
				log.Fatal("html is malformed")
			}
			return
		default:
			book := parseBook(s, pd.Url)
			o <- book
		}
	}
}

func parseBook(s *tokenStream, url string) domain.Book {
	attr, err := s.skipToBeginning("p", map[string]string{})
	if err != nil {
		log.Fatal("html is malformed")
	}

	rateings := parseRatings(attr["class"])

	attr, err = s.skipToBeginning("a", map[string]string{})
	if err != nil {
		log.Fatal("html is malformed")
	}

	name := attr["title"]
	href := attr["href"]

	if name == "" || href == "" {
		log.Fatal("html is malformed")
	}

	details := resolveUrl(url, href)

	_, err = s.skipToBeginning("p", map[string]string{"class": "price_color"})
	if err != nil {
		log.Fatal("html is malformed")
	}

	if s.Next() != html.TextToken {
		log.Fatal("html is malformed")
	}
	price, err := domain.NewMoney(string(s.z.Text()))
	if err != nil {
		log.Fatal("html is malformed")
	}

	attr, err = s.skipToBeginning("p", map[string]string{})
	if err != nil {
		log.Fatal("html is malformed")
	}

	inStock := parseInStock(attr["class"])
	defer s.skipToEnd("li")

	return domain.Book{
		Name:       name,
		Rating:     rateings,
		Price:      price,
		InStock:    inStock,
		DetailPage: details,
	}
}

func parseRatings(ratingClass string) domain.Rating {
	a, f := strings.CutPrefix(ratingClass, "star-rating ")
	if !f {
		log.Fatal("html is malformed")
	}

	r := uint8(0)
	switch a {
	case "One":
		r = 1
	case "Two":
		r = 2
	case "Three":
		r = 3
	case "Four":
		r = 4
	case "Five":
		r = 5
	}

	rating, _ := domain.NewRateing(r)
	return rating
}

func parseInStock(class string) bool {
	a, f := strings.CutSuffix(class, " availability")
	if !f {
		log.Fatal("html is malformed")
	}

	return a == "instock"
}
