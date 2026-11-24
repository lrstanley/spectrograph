// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"errors"
	"fmt"
	"testing"

	"github.com/diamondburned/arikawa/v3/utils/ws"
)

func TestIsUnknownEvent(t *testing.T) {
	t.Parallel()

	unknown := ws.UnknownEventError{Op: 0, Type: "VOICE_CHANNEL_START_TIME_UPDATE"}

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "value", err: unknown, want: true},
		{name: "pointer", err: &unknown, want: true},
		{name: "wrapped", err: fmt.Errorf("decode: %w", unknown), want: true},
		{name: "other", err: errors.New("connection reset"), want: false},
		{name: "nil", err: nil, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := isUnknownEvent(tt.err); got != tt.want {
				t.Errorf("isUnknownEvent(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
