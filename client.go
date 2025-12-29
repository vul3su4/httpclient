package httpclient

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Config struct{
	Timeout time.Duration
	ProxyURL string
	BaseHeaders map[string]string
	MaxIdleConns int
	MaxIdleConnsPerHost int
	IdleConnTimeout time.Duration
}

type Client struct{
	http *http.Client
	baseHeaders map[string]string
}


func New(cfg Config) (*Client, error) {
	if cfg.Timeout == 0 {
		cfg.Timeout = 15 * time.Second
	}
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 1000
	}
	if cfg.MaxIdleConnsPerHost == 0 {
		cfg.MaxIdleConnsPerHost = 100
	}
	if cfg.IdleConnTimeout == 0 {
		cfg.IdleConnTimeout = 90 * time.Second
	}

	jar, _ := cookiejar.New(nil)

	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        cfg.MaxIdleConns,
		MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
		IdleConnTimeout:     cfg.IdleConnTimeout,
		TLSHandshakeTimeout: 10 * time.Second,
		Proxy:              http.ProxyFromEnvironment,
	}

	if cfg.ProxyURL != "" {
		p, err := url.Parse(cfg.ProxyURL)
		if err != nil {
			return nil, err
		}
		tr.Proxy = http.ProxyURL(p)
	}

	c := &Client{
		http: &http.Client{
			Timeout:   cfg.Timeout,
			Transport: tr,
			Jar:       jar,
		},
		baseHeaders: map[string]string{
			"User-Agent": "Mozilla/5.0",
		},
	}

	for k, v := range cfg.BaseHeaders {
		c.baseHeaders[k] = v
	}

	return c, nil
}