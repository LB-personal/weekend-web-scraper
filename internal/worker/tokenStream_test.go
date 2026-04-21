package worker

import (
	"io"
	"reflect"
	"testing"

	"golang.org/x/net/html"
)

func Test_tokenStream_Next(t *testing.T) {
	tests := []struct {
		name string
		s    *tokenStream
		want []html.TokenType
		err  *error
	}{
		{"plain text", newTokenStream([]byte("hello world")), []html.TokenType{html.TextToken, html.ErrorToken}, &io.EOF},
		{"white spaces", newTokenStream([]byte("  \t")), []html.TokenType{html.TextToken, html.ErrorToken}, &io.EOF},
		{"numbers and punctuations", newTokenStream([]byte("  \t")), []html.TokenType{html.TextToken, html.ErrorToken}, &io.EOF},
		{"ignore cases", newTokenStream([]byte("<!DOCTYPE html><!-- this is a comment -->")), []html.TokenType{html.ErrorToken}, &io.EOF},
		{"wellformed html", newTokenStream([]byte("<div><p>text</p></div>")), []html.TokenType{html.StartTagToken, html.StartTagToken, html.TextToken, html.EndTagToken, html.EndTagToken, html.ErrorToken}, &io.EOF},
		{"self closing tags", newTokenStream([]byte("<br/><img/><input/>")), []html.TokenType{html.SelfClosingTagToken, html.SelfClosingTagToken, html.SelfClosingTagToken, html.ErrorToken}, &io.EOF},
		{"single attribute", newTokenStream([]byte("<a href=\"url\">")), []html.TokenType{html.StartTagToken, html.ErrorToken}, &io.EOF},
		{"multiple attributes", newTokenStream([]byte("<input type=\"text\" id=\"name\" required>")), []html.TokenType{html.StartTagToken, html.ErrorToken}, &io.EOF},
		{"boolean attributes", newTokenStream([]byte("<input disabled>")), []html.TokenType{html.StartTagToken, html.ErrorToken}, &io.EOF},
		{"single-quoted values attributes", newTokenStream([]byte("<div class='container'>")), []html.TokenType{html.StartTagToken, html.ErrorToken}, &io.EOF},
		{"unclosed tags", newTokenStream([]byte("<div><p>text</div>")), []html.TokenType{html.StartTagToken, html.StartTagToken, html.TextToken, html.EndTagToken, html.ErrorToken}, &io.EOF},
		{"misnested tags", newTokenStream([]byte("<b><i>text</b></i>")), []html.TokenType{html.StartTagToken, html.StartTagToken, html.TextToken, html.EndTagToken, html.EndTagToken, html.ErrorToken}, &io.EOF},
		{"missing closing bracket", newTokenStream([]byte("<div class=\"foo\"")), []html.TokenType{html.ErrorToken}, &io.EOF},
		{"custom elements", newTokenStream([]byte("<my-component>")), []html.TokenType{html.StartTagToken, html.ErrorToken}, &io.EOF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, w := range tt.want {
				got := tt.s.Next()
				if got != w {
					t.Errorf("tokenStream.Next() expectec %v, actual %v", tt.want, got)
				}

				if got == html.ErrorToken && tt.s.z.Err() != *tt.err {
					t.Errorf("tokenStream.Next() expected %v error, actual %v", tt.want, got)
				}
			}
		})
	}
}

func Test_tokenStream_Tag(t *testing.T) {
	tests := []struct {
		name  string
		s     *tokenStream
		want  string
		want1 map[string]string
	}{
		{"single attribute", newTokenStream([]byte("<a href=\"url\">")), "a", map[string]string{"href": "url"}},
		{"multiple attributes", newTokenStream([]byte("<input type=\"text\" id=\"name\" required>")), "input", map[string]string{"type": "text", "id": "name", "required": ""}},
		{"boolean attributes", newTokenStream([]byte("<input disabled>")), "input", map[string]string{"disabled": ""}},
		{"single-quoted values attributes", newTokenStream([]byte("<div class='container'>")), "div", map[string]string{"class": "container"}},
		{"empty value", newTokenStream([]byte("<div class=\"\">")), "div", map[string]string{"class": ""}},
		{"duplicate attributes", newTokenStream([]byte("<div id=\"a\" id=\"b\">")), "div", map[string]string{"id": "a"}},
		{"attribute with no space", newTokenStream([]byte("<div id=\"a\"class=\"b\">")), "div", map[string]string{"id": "a", "class": "b"}},
		{"extra whitespace in tag", newTokenStream([]byte("<div   class = \"foo\"  >")), "div", map[string]string{"class": "foo"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Next()
			got, got1 := tt.s.Tag()
			if got != tt.want {
				t.Errorf("tokenStream.Tag() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("tokenStream.Tag() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_isMapSubset(t *testing.T) {
	type args struct {
		m   map[string]string
		sub map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"equal", args{map[string]string{"a": "1", "b": "2"}, map[string]string{"a": "1", "b": "2"}}, true},
		{"true subset", args{map[string]string{"a": "1", "b": "2"}, map[string]string{"b": "2"}}, true},
		{"nil set to a set", args{map[string]string{"a": "1", "b": "2"}, map[string]string{}}, true},
		{"nil set to nil set", args{map[string]string{}, map[string]string{}}, true},
		{"complete mismatch", args{map[string]string{"a": "1", "b": "2"}, map[string]string{"c": "3", "d": "4"}}, false},
		{"partial mismatch", args{map[string]string{"a": "1", "b": "2"}, map[string]string{"c": "3", "b": "2"}}, false},
		{"superset", args{map[string]string{"b": "2"}, map[string]string{"a": "1", "b": "2"}}, false},
		{"a set to nil set", args{map[string]string{}, map[string]string{"a": "1", "b": "2"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMapSubset(tt.args.m, tt.args.sub); got != tt.want {
				t.Errorf("isMapSubset() = %v, want %v", got, tt.want)
			}
		})
	}
}
