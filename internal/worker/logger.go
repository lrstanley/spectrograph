// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/diamondburned/arikawa/v3/utils/ws"
)

var setLoggerOnce sync.Once

// setLogger routes arikawa's websocket logging (package-level hooks, which
// default to the standard library logger) through slog. Only the first call
// has an effect, as the hooks are global.
func setLogger(logger *slog.Logger) {
	setLoggerOnce.Do(func() {
		ws.WSError = func(err error) {
			logger.Error("gateway error", "error", err)
		}
		ws.WSDebug = func(v ...any) {
			logger.Debug(strings.TrimSpace(fmt.Sprintln(v...))) //nolint:sloglint // dynamic messages from arikawa.
		}
	})
}
