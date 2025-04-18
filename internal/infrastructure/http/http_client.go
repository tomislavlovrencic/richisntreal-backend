package http

import (
	"net/http"
	"sync"
)

var (
	httpclientSingleton *http.Client
	httpclientOnce      = &sync.Once{}
)

func NewHTTPClient() *http.Client {
	httpclientOnce.Do(func() {
		httpclientSingleton = &http.Client{}
	})

	return httpclientSingleton
}
