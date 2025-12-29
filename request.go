package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Options struct {
	Headers map[string]string
	Query   map[string]string
	Cookies []*http.Cookie
}

func (c *Client) Do(ctx context.Context, method, rawURL string, body io.Reader, opt Options) (*http.Response, error) {
	finalURL, err := applyQuery(rawURL, opt.Query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, finalURL, body)
	if err != nil {
		return nil, err
	}

	// base headers
	for k, v := range c.baseHeaders {
		req.Header.Set(k, v)
	}
	// per-request headers
	for k, v := range opt.Headers {
		req.Header.Set(k, v)
	}
	// cookies (optional)
	for _, ck := range opt.Cookies {
		req.AddCookie(ck)
	}

	return c.http.Do(req)
}

func (c *Client) Get(ctx context.Context, rawURL string, opt Options) (*http.Response, error) {
	return c.Do(ctx, http.MethodGet, rawURL, nil, opt)
}

func (c *Client) Post(ctx context.Context, rawURL string, body io.Reader, opt Options) (*http.Response, error) {
	return c.Do(ctx, http.MethodPost, rawURL, body, opt)
}

// PostJSON:  marshal + set Content-Type + decode（out can be nil）
func (c *Client) PostJSON(ctx context.Context, rawURL string, payload any, out any, opt Options) (*http.Response, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	if opt.Headers == nil {
		opt.Headers = map[string]string{}
	}
	if opt.Headers["Content-Type"] == "" {
		opt.Headers["Content-Type"] = "application/json"
	}

	resp, err := c.Do(ctx, http.MethodPost, rawURL, bytes.NewReader(b), opt)
	if err != nil {
		return nil, err
	}

	if out != nil {
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return nil, err
		}
	}
	return resp, nil
}

func applyQuery(rawURL string, query map[string]string) (string, error) {
	if len(query) == 0 {
		return rawURL, nil
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
