// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package generate

//go:generate sh -c "protoc --proto_path=. --go_out=paths=source_relative:../ *.proto"
//go:generate protoc-go-inject-tag -input=../*.pb.go