// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import "regexp"

const (
	DefaultMaxClones   = 20
	DefaultMaxChannels = 10
)

var DefaultChannelMatch = regexp.MustCompile(`^.* \+$`)
