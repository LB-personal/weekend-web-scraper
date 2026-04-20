package worker

import (
	"bytes"
	"io"
	"log"

	"golang.org/x/net/html"
)

type tokenStream struct {
	z *html.Tokenizer
}

func newTokenStream(b []byte) *tokenStream {
	return &tokenStream{
		z: html.NewTokenizer(bytes.NewReader(b)),
	}
}

func (s *tokenStream) Next() html.TokenType {
	for {
		tt := s.z.Next()
		switch tt {
		case html.TextToken, html.CommentToken, html.DoctypeToken:
			continue
		case html.ErrorToken:
			err := s.z.Err()
			if err != io.EOF {
				log.Fatal("scraped site is broken", s.z.Err())
			}
		default:
			return tt
		}
	}
}

func (s *tokenStream) Tag() (string, map[string]string) {
	name, hasAttr := s.z.TagName()
	attrs := map[string]string{}
	for hasAttr {
		k, v, more := s.z.TagAttr()
		attrs[string(k)] = string(v)
		hasAttr = more
	}
	return string(name), attrs
}

func (s *tokenStream) skipToBeginning(tag string, attr map[string]string) (map[string]string, error) {
	for {
		switch s.Next() {
		case html.StartTagToken:
			n, a := s.Tag()
			if n == tag && isMapSubset(a, attr) {
				return a, nil
			}
		case html.ErrorToken:
			return nil, io.EOF
		}
	}
}

func (s *tokenStream) skipToEnd(tag string) error {
	for {
		switch s.Next() {
		case html.EndTagToken:
			n, _ := s.Tag()
			if n == tag {
				return nil
			}
		case html.ErrorToken:
			return io.EOF
		}
	}
}

func isMapSubset[K, V comparable](m, sub map[K]V) bool {
	if len(sub) > len(m) {
		return false
	}
	for k, vsub := range sub {
		if vm, found := m[k]; !found || vm != vsub {
			return false
		}
	}
	return true
}
