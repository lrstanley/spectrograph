// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import (
	"testing"

	"github.com/alecthomas/kong"
)

// parseFlags parses args the same way clix does (the flags embedded into the
// root struct), so the kong hooks on embedded structs are exercised.
func parseFlags(t *testing.T, args ...string) (*Flags, error) {
	t.Helper()

	root := struct {
		Flags *Flags `embed:""`
	}{Flags: &Flags{}}

	parser, err := kong.New(&root, kong.Exit(func(int) { t.Fatal("unexpected exit") }))
	if err != nil {
		t.Fatal(err)
	}
	_, err = parser.Parse(append([]string{
		"--http.validation-key=a",
		"--http.encryption-key=b",
		"--discord.client-id=c",
		"--discord.client-secret=d",
	}, args...))
	return root.Flags, err
}

func TestFlagsURLs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		args        []string
		wantDB      string
		wantBaseURL string
		wantErr     bool
	}{
		{
			name:        "defaults",
			args:        []string{"--database.url=postgres://user:pass@host:5432/db?sslmode=require&channel_binding=require"},
			wantDB:      "postgres://user:pass@host:5432/db?sslmode=require&channel_binding=require",
			wantBaseURL: "http://localhost:8080",
		},
		{
			name: "database-overrides",
			args: []string{
				"--database.url=postgres://user:pass@host:5432/db?sslmode=require",
				"--database.username=other",
				"--database.database=otherdb",
			},
			wantDB:      "postgres://other:pass@host:5432/otherdb?sslmode=require",
			wantBaseURL: "http://localhost:8080",
		},
		{
			name: "database-credentials",
			args: []string{
				"--database.url=postgres://host/db",
				"--database.username=u",
				"--database.password=p@ss",
				"--http.base-url=https://example.com/",
			},
			wantDB:      "postgres://u:p%40ss@host/db",
			wantBaseURL: "https://example.com/",
		},
		{
			name:    "relative-base-url",
			args:    []string{"--database.url=postgres://host/db", "--http.base-url=example.com"},
			wantErr: true,
		},
		{
			name:    "non-http-base-url",
			args:    []string{"--database.url=postgres://host/db", "--http.base-url=ftp://example.com"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			flags, err := parseFlags(t, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if got := flags.Database.URL.String(); got != tt.wantDB {
				t.Errorf("database url = %q, want %q", got, tt.wantDB)
			}
			if got := flags.HTTP.BaseURL.String(); got != tt.wantBaseURL {
				t.Errorf("base url = %q, want %q", got, tt.wantBaseURL)
			}
		})
	}
}
