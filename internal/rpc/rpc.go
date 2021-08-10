// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package rpc

import (
	http "net/http"
	"os"
	strconv "strconv"
	"time"

	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/gojek/heimdall/v7/plugins"
	"github.com/twitchtv/twirp"
)

//go:generate protoc --go_out=paths=source_relative:. --twirp_out=paths=source_relative:. service.proto

// PathPrefixWorker is the prefix used for calls to the rpc server. This does not
// include any other prefixes that may be needed to mount the server on the
// http server mux.
const PathPrefixWorker = "/api/rpc/worker"

func NewWorkerClient(baseURL, secretKey, version string, shardID int, timeout time.Duration, maxRetries int) Worker {
	initalTimeout := 1 * time.Second         // Inital timeout
	maxTimeout := 30 * time.Second           // Max time out
	exponentFactor := 2.0                    // Multiplier
	maximumJitterInterval := 1 * time.Second // Max jitter interval.

	backoff := heimdall.NewExponentialBackoff(initalTimeout, maxTimeout, exponentFactor, maximumJitterInterval)
	retrier := heimdall.NewRetrier(backoff)

	client := httpclient.NewClient(
		httpclient.WithHTTPClient(&apiClient{client: http.DefaultClient, key: secretKey, version: version, shardId: shardID}),
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(maxRetries),
	)
	// TODO: swap out with a leveled implementation:
	// https://github.com/gojek/heimdall/blob/master/plugins/request_logger.go
	client.AddPlugin(plugins.NewRequestLogger(os.Stdout, nil))

	return NewWorkerProtobufClient(baseURL, client, twirp.WithClientPathPrefix(PathPrefixWorker))
}

type apiClient struct {
	client  *http.Client
	version string
	key     string
	shardId int
}

func (c *apiClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Api-Version", c.version)
	req.Header.Set("X-Api-Key", c.key)
	req.Header.Set("X-Shard-Id", strconv.Itoa(c.shardId))
	return c.client.Do(req)
}
