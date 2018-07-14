package ihttp_test

import (
	"net/http"
	"testing"

	"bytes"
	"io/ioutil"

	"time"

	"sync"

	"github.com/hafidsousa/lib-cryptocurrency/ihttp"
	"github.com/stretchr/testify/assert"
)

type mockTransport struct{}

func newMockTransport() http.RoundTripper {
	return &mockTransport{}
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString("Test Body")),
	}, nil
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		opt         ihttp.Options
		attempts    int
		minDuration time.Duration
		maxDuration time.Duration
	}{
		{
			name:        "shouldAllowUnlimitedRequests",
			opt:         ihttp.Options{},
			attempts:    100,
			minDuration: 1 * time.Nanosecond,
			maxDuration: 1 * time.Millisecond,
		},
		{
			name:        "shouldLimitRequests",
			opt:         ihttp.Options{Limit: 100},
			attempts:    400,
			minDuration: 3 * time.Second,
			maxDuration: 4 * time.Second,
		},
	}

	var wg sync.WaitGroup

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := &http.Client{Transport: newMockTransport()}
			ic := ihttp.NewClient(c, test.opt)
			start := time.Now()
			wg.Add(test.attempts)

			for i := 1; i <= test.attempts; i++ {
				go job(t, ic, &wg)
			}
			// Wait for all HTTP fetches to complete.
			wg.Wait()
			assert.Condition(t, func() bool {
				elapsed := time.Since(start)
				return elapsed > test.minDuration && elapsed < test.maxDuration
			})
		})
	}
}

func job(t *testing.T, ic ihttp.Client, wg *sync.WaitGroup) {
	req, _ := http.NewRequest("GET", "http://localhost:8080", nil)
	resp, err := ic.DoRequest(req)
	if err != nil {
		panic(err)
	}
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	defer wg.Done()
}
