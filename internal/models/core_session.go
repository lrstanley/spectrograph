// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import (
	"time"
)

// Session represents a scs session needed for encoding/decoding to/from a database.
type Session struct {
	Token   string    `bson:"token"`
	Data    []byte    `bson:"data"`
	Expires time.Time `bson:"expires"`
}

const (
	SessionUserIDKey = "user_id"
	SessionAdminKey  = "isadmin"
)
