package apiclient

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/gojek/heimdall/v7/plugins"
)

func New(baseURL string, timeout time.Duration, maxRetries int, headers map[string]string) *httpclient.Client {
	initalTimeout := 1 * time.Second         // Inital timeout.
	maxTimeout := 30 * time.Second           // Max time out.
	exponentFactor := 2.0                    // Multiplier.
	maximumJitterInterval := 1 * time.Second // Max jitter interval.

	backoff := heimdall.NewExponentialBackoff(initalTimeout, maxTimeout, exponentFactor, maximumJitterInterval)
	retrier := heimdall.NewRetrier(backoff)

	client := httpclient.NewClient(
		httpclient.WithHTTPClient(&apiClient{client: http.DefaultClient, baseURL: baseURL, headers: headers}),
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(maxRetries),
	)

	// TODO: swap out with a leveled implementation:
	// https://github.com/gojek/heimdall/blob/master/plugins/request_logger.go
	client.AddPlugin(plugins.NewRequestLogger(os.Stdout, nil))

	return client
}

type apiClient struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

func (c *apiClient) Do(req *http.Request) (*http.Response, error) {
	var err error

	// Merge baseURL + requested url.
	req.URL, err = url.Parse(strings.TrimRight(c.baseURL, "/") + req.URL.String())
	if err != nil {
		return nil, err
	}

	if c.headers != nil {
		if req.Header == nil {
			req.Header = http.Header{}
		}
		for k, v := range c.headers {
			req.Header.Set(k, v)
		}
	}
	return c.client.Do(req)
}
