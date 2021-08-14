// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import (
	"context"
)

type ServerService interface {
	Upsert(ctx context.Context, r *Server) error
	Get(ctx context.Context, id string) (*Server, error)
	List(ctx context.Context) ([]*Server, error)
}

func (s *ServerDiscordData) ParsePermissions() DiscordPermissions {
	return DiscordPermissions(s.Permissions)
}
