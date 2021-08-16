// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ServerService interface {
	Upsert(ctx context.Context, r *Server) error
	Get(ctx context.Context, id string) (*Server, error)
	GetByDiscordID(ctx context.Context, id string) (*Server, error)
	List(ctx context.Context) ([]*Server, error)
}

func (s *Server) Validate() error {
	if s.Created == nil {
		s.Created = timestamppb.Now()
	}

	if s.Updated == nil {
		s.Updated = s.Created
	}

	return errUseBuiltinValidator
}

func (s *ServerDiscordData) ParsePermissions() DiscordPermissions {
	return DiscordPermissions(s.Permissions)
}
