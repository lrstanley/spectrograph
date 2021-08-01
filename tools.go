// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

// This file is used to tell go mod to keep track of the exact cli tool versions
// that are used for the project. This file itself doesn't install these tools,
// only ensure the correct versions are pinned.
//
// See relevant Makefile's for "go install" commands.

package tools

// +build tools
import (
	_ "github.com/GeertJohan/go.rice/rice"
	_ "github.com/twitchtv/twirp/protoc-gen-twirp"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
