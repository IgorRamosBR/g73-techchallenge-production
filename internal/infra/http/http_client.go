package http

import (
	"bytes"
	httpClient "net/http"
	"time"
)

type HttpClient interface {
	DoGet(url string) (*httpClient.Response, error)
	DoPut(url string, body []byte) (*httpClient.Response, error)
}

type client struct {
	client *httpClient.Client
}

func NewHttpClient(timeout time.Duration) HttpClient {
	return client{
		client: &httpClient.Client{
			Timeout: timeout,
		},
	}
}

func (c client) DoPut(url string, body []byte) (*httpClient.Response, error) {
	req, err := httpClient.NewRequest(httpClient.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}

func (c client) DoGet(url string) (*httpClient.Response, error) {
	return httpClient.Get(url)
}
