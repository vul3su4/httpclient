package httpclient

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
	"fmt"
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

func (c *Client) SetCookie(rawURL, name, value string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	c.http.Jar.SetCookies(u, []*http.Cookie{
		{Name: name, Value: value, Path: "/"},
	})
	return nil
}

func (c *Client) DeleteCookie(rawURL, name string) error {
    u, err := url.Parse(rawURL)
    if err != nil {
        return err
    }

    // MaxAge < 0 represents deleting the cookie
    c.http.Jar.SetCookies(u, []*http.Cookie{
        {
            Name:    name,
            Value:   "",
            Path:    "/",
            MaxAge:  -1,
            Expires: time.Unix(0, 0),
        },
    })
    return nil
}

func (c *Client) DumpCookies(rawURL string) error {
    u, err := url.Parse(rawURL)
    if err != nil {
        return err
    }

    cookies := c.http.Jar.Cookies(u)
    if len(cookies) == 0 {
        fmt.Println("[DumpCookies] (no cookies)")
        return nil
    }

    fmt.Println("[DumpCookies]")
    for _, ck := range cookies {
        fmt.Printf("  %s=%s\n", ck.Name, ck.Value)
    }
    return nil
}

func (c *Client) GetCookies(rawURL string) ([]*http.Cookie, error) {
    u, err := url.Parse(rawURL)
    if err != nil {
        return nil, err
    }
    return c.http.Jar.Cookies(u), nil
}