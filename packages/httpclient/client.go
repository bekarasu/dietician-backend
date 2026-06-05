package httpclient

import (
	"net/http"
	"time"
)

func New(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

func Default() *http.Client {
	return New(10 * time.Second)
}
