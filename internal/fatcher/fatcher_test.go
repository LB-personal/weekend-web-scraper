package fatcher

import (
	"context"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
	"testing"

	"github.com/LB-personal/weekend-web-scraper/internal/domain"
)

type cannedResponse struct {
	resp     *http.Response
	err      error
	hitCount int
}

type mockHttpGetter struct {
	cannedResponses map[string]*cannedResponse
}

func (m *mockHttpGetter) Get(url string) (resp *http.Response, err error) {
	if val, ok := m.cannedResponses[url]; !ok {
		log.Fatal("mock HttpGeter got a url request it was not setup for")
		return nil, nil // silence the compiler
	} else {
		val.hitCount++
		return val.resp, val.err
	}
}

type testScenario struct {
	name     string
	in       chan string
	f        *fatcher
	inputs   []string
	expected []*domain.PageData
}

func newTestScenario(name string, inputs []string, expected []*domain.PageData, cr map[string]*cannedResponse) testScenario {
	in := make(chan string)
	return testScenario{
		name: name,
		in:   in,
		f: &fatcher{
			client: &mockHttpGetter{
				cannedResponses: cr,
			},
			urls:       in,
			bufferSize: 1,
		},
		inputs:   inputs,
		expected: expected,
	}
}

func Test_fatcher_Fatch(t *testing.T) {
	tests := []testScenario{
		newTestScenario(
			"fetch one url",
			[]string{"mock"},
			[]*domain.PageData{
				{
					Content: []byte{},
					Url:     "mock",
				},
			},
			map[string]*cannedResponse{
				"mock": {
					resp: &http.Response{
						StatusCode: 200,
						Status:     "Ok",
						Body:       io.NopCloser(strings.NewReader("")),
					},
					err:      nil,
					hitCount: 0,
				},
			},
		),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.f.Fatch(context.Background())
			for i, v := range tt.inputs {
				tt.in <- v
				pd := <-out
				if !slices.Equal(pd.Content, tt.expected[i].Content) {
					t.Errorf("%s: iteration %d expected %+q content, actual %+q", tt.name, i, tt.expected[i].Content, pd.Content)
				}
				if v != pd.Url {
					t.Errorf("%s: iteration %d returned a url that is diffrent then what sent to it", tt.name, i)
				}
			}

			close(tt.in)
		})
	}
}

// need to test negative tests but plan to change error handling in the future
