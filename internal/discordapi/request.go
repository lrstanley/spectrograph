// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

// Package discordapi provides a minimal Discord REST API client, used with
// user-provided OAuth2 bearer tokens (the bot itself uses the worker's client).
package discordapi

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"time"
)

// APIBase is the versioned Discord API base URL.
const APIBase = "https://discord.com/api/v10"

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func request[T any](ctx context.Context, token, method, url string, result *T) error {
	req, err := http.NewRequestWithContext(ctx, method, url, http.NoBody)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from %s: %d", url, resp.StatusCode)
	}

	if err = json.UnmarshalRead(resp.Body, result); err != nil {
		return fmt.Errorf("decoding response from %s: %w", url, err)
	}
	return nil
}
