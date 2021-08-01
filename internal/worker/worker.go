// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

// PathPrefix is the prefix used for calls to the rpc server. This does not
// include any other prefixes that may be needed to mount the server on the
// http server mux.
const PathPrefix = "/api/rpc/worker"

//go:generate protoc --go_out=paths=source_relative:. --twirp_out=paths=source_relative:. service.proto
