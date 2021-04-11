// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

// ContextKey is a type that prevent key collisions in contexts. When the type
// is different, even if the key name is the same, it will never overlap with
// another package.
type ContextKey string

const (
	ContextDebug ContextKey = "debug"
)
