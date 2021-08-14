// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package rpc

import (
	http "net/http"
	"os"
	"time"

	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/gojek/heimdall/v7/plugins"
	"github.com/twitchtv/twirp"
)

//go:generate sh -c "cd ../../;protoc --proto_path=. --go_out=paths=source_relative:. --twirp_out=paths=source_relative:. internal/rpc/*.proto"
//go:generate protoc-go-inject-tag -input=*.pb.go

// PathPrefix is the prefix used for calls to the rpc server. This is the
// absolute path on the http server (baseURL).
const PathPrefix = "/api/rpc"

func NewWorkerClient(baseURL string, client *httpclient.Client) Worker {
	return NewWorkerProtobufClient(baseURL, client, twirp.WithClientPathPrefix(PathPrefix))
}

func NewHTTPClient(timeout time.Duration, headers map[string]string) *httpclient.Client {
	initialTimeout := 500 * time.Millisecond // Initial timeout.
	maxTimeout := 15 * time.Second           // Max timeout.
	exponentFactor := 2.0                    // Multiplier.
	maximumJitterInterval := 1 * time.Second // Max jitter interval.

	backoff := heimdall.NewExponentialBackoff(initialTimeout, maxTimeout, exponentFactor, maximumJitterInterval)
	client := httpclient.NewClient(
		httpclient.WithHTTPClient(&rpcClient{client: http.DefaultClient, headers: headers}),
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetrier(heimdall.NewRetrier(backoff)),
		httpclient.WithRetryCount(5),
	)

	// TODO: swap out with a leveled implementation:
	// https://github.com/gojek/heimdall/blob/master/plugins/request_logger.go
	client.AddPlugin(plugins.NewRequestLogger(os.Stdout, nil))

	return client
}

type rpcClient struct {
	client  *http.Client
	headers map[string]string
}

func (c *rpcClient) Do(req *http.Request) (*http.Response, error) {
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
