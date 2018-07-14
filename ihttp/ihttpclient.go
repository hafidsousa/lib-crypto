// Package ihttp is an internal http.Client that encapsulates Golang `http.Client` functions
package ihttp

import (
	"net/http"

	"errors"

	"go.uber.org/ratelimit"
)

type Client interface {
	DoRequest(req *http.Request) (*http.Response, error)
}

type DefaultClient struct {
	client  *http.Client
	limiter ratelimit.Limiter
}

type Options struct {
	// Limit represents the rate-limiting of requests per seconds
	// A Limit of zero means no rate-limiting
	Limit int
}

// NewClient creates a new HttpClient
func NewClient(c *http.Client, opt Options) Client {
	rl := (ratelimit.Limiter)(nil)
	if opt.Limit != 0 {
		rl = ratelimit.New(opt.Limit)
	}
	return &DefaultClient{
		client:  c,
		limiter: rl,
	}
}

// DoRequest encapsulates Go `http.Client.Do`
func (c *DefaultClient) DoRequest(req *http.Request) (*http.Response, error) {
	// Apply rate limiting
	if c.limiter != nil {
		c.limiter.Take()
	}

	resp, err := c.client.Do(req)

	if err != nil {
		err = errors.New("error making http request")
	}
	defer resp.Body.Close()

	return resp, err
}
