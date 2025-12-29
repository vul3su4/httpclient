package httpclient


import (
	"io"
	"net/http"
)

func ReadBody(resp *http.Response, maxBytes int) (string, error) {
	defer resp.Body.Close()
	limited := io.LimitReader(resp.Body, int64(maxBytes))
	b, err := io.ReadAll(limited)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
