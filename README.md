# httpclient

HTTP client library for Go, designed to simplify HTTP requests with built-in cookie management, proxy support, and connection pooling.

## Features

-  **Easy to use**: Simple API for making HTTP requests
-  **Cookie Management**: Built-in cookie jar with helper methods
-  **Configurable**: Customizable timeout, proxy, headers, and connection pool settings
-  **JSON Support**: Convenient `PostJSON` method for JSON requests/responses
-  **Connection Pooling**: Optimized connection reuse with configurable limits
-  **Context Support**: Full support for Go contexts for cancellation and timeouts

## Installation

```bash
go get github.com/vul3su4/httpclient
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "github.com/vul3su4/httpclient"
    "time"
)

func main() {
    // Create a new client with default settings
    client, err := httpclient.New(httpclient.Config{})
    if err != nil {
        panic(err)
    }

    ctx := context.Background()

    // Simple GET request
    resp, err := client.Get(ctx, "https://api.example.com/data", httpclient.Options{})
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, err := httpclient.ReadBody(resp, 1024*1024) // Read up to 1MB
    if err != nil {
        panic(err)
    }
    fmt.Println(body)
}
```

### Advanced Configuration

```go
// Create a client with custom configuration
client, err := httpclient.New(httpclient.Config{
    Timeout:            30 * time.Second,
    ProxyURL:          "http://proxy.example.com:8080", // Optional
    MaxIdleConns:      1000,
    MaxIdleConnsPerHost: 100,
    IdleConnTimeout:   90 * time.Second,
    BaseHeaders: map[string]string{
        "User-Agent": "MyApp/1.0",
        "Accept":     "application/json",
    },
})
```

### GET Request with Query Parameters

```go
ctx := context.Background()

resp, err := client.Get(ctx, "https://api.example.com/search", httpclient.Options{
    Query: map[string]string{
        "q":     "golang",
        "limit": "10",
        "page":  "1",
    },
    Headers: map[string]string{
        "Authorization": "Bearer your-token",
    },
})
if err != nil {
    panic(err)
}
defer resp.Body.Close()

body, _ := httpclient.ReadBody(resp, 1024*1024)
fmt.Println(body)
```

### POST Request with JSON

```go
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token string `json:"token"`
    User  struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
    } `json:"user"`
}

ctx := context.Background()

loginReq := LoginRequest{
    Username: "user@example.com",
    Password: "secret123",
}

var loginResp LoginResponse
resp, err := client.PostJSON(ctx, "https://api.example.com/login", loginReq, &loginResp, httpclient.Options{
    Headers: map[string]string{
        "X-API-Key": "your-api-key",
    },
})
if err != nil {
    panic(err)
}
defer resp.Body.Close()

fmt.Printf("Token: %s\n", loginResp.Token)
fmt.Printf("User: %s (ID: %d)\n", loginResp.User.Name, loginResp.User.ID)
```

### POST Request with Raw Body

```go
ctx := context.Background()

bodyData := []byte("name=John&email=john@example.com")
body := bytes.NewReader(bodyData)

resp, err := client.Post(ctx, "https://api.example.com/users", body, httpclient.Options{
    Headers: map[string]string{
        "Content-Type": "application/x-www-form-urlencoded",
    },
})
if err != nil {
    panic(err)
}
defer resp.Body.Close()

result, _ := httpclient.ReadBody(resp, 1024*1024)
fmt.Println(result)
```

### Cookie Management

```go
// Set a cookie
err := client.SetCookie("https://example.com", "session_id", "abc123xyz")
if err != nil {
    panic(err)
}

// Get all cookies for a URL
cookies, err := client.GetCookies("https://example.com")
if err != nil {
    panic(err)
}
for _, cookie := range cookies {
    fmt.Printf("Cookie: %s=%s\n", cookie.Name, cookie.Value)
}

// Dump cookies to console (for debugging)
client.DumpCookies("https://example.com")

// Delete a cookie
err = client.DeleteCookie("https://example.com", "session_id")
if err != nil {
    panic(err)
}
```

### Using Context for Timeout

```go
// Create a context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := client.Get(ctx, "https://api.example.com/slow-endpoint", httpclient.Options{})
if err != nil {
    if err == context.DeadlineExceeded {
        fmt.Println("Request timed out")
    } else {
        panic(err)
    }
}
defer resp.Body.Close()
```

### Custom Request with Do Method

```go
ctx := context.Background()

// PUT request
data := bytes.NewReader([]byte("update data"))
resp, err := client.Do(ctx, http.MethodPut, "https://api.example.com/resource/123", data, httpclient.Options{
    Headers: map[string]string{
        "Content-Type": "application/json",
    },
})
if err != nil {
    panic(err)
}
defer resp.Body.Close()

// DELETE request
resp, err = client.Do(ctx, http.MethodDelete, "https://api.example.com/resource/123", nil, httpclient.Options{})
if err != nil {
    panic(err)
}
defer resp.Body.Close()
```


### Config

Configuration options for creating a new client:

```go
type Config struct {
    Timeout            time.Duration  // Request timeout (default: 15s)
    ProxyURL           string         // Proxy URL (optional)
    BaseHeaders        map[string]string // Default headers for all requests
    MaxIdleConns       int            // Max idle connections (default: 1000)
    MaxIdleConnsPerHost int           // Max idle connections per host (default: 100)
    IdleConnTimeout    time.Duration  // Idle connection timeout (default: 90s)
}
```

### Options

Request options for individual requests:

```go
type Options struct {
    Headers map[string]string  // Request headers
    Query   map[string]string  // Query parameters
    Cookies []*http.Cookie     // Additional cookies
}
```

### Methods

- `New(cfg Config) (*Client, error)` - Create a new HTTP client
- `Get(ctx, url, opt Options) (*http.Response, error)` - Make a GET request
- `Post(ctx, url, body, opt Options) (*http.Response, error)` - Make a POST request
- `PostJSON(ctx, url, payload, out, opt Options) (*http.Response, error)` - POST JSON and decode response
- `Do(ctx, method, url, body, opt Options) (*http.Response, error)` - Make a custom HTTP request
- `SetCookie(url, name, value string) error` - Set a cookie
- `DeleteCookie(url, name string) error` - Delete a cookie
- `GetCookies(url string) ([]*http.Cookie, error)` - Get all cookies for a URL
- `DumpCookies(url string) error` - Print all cookies to console

### Helper Functions

- `ReadBody(resp *http.Response, maxBytes int) (string, error)` - Read response body with size limit



