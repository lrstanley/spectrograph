// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/httputil/httpdriver"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/models"
)

// offlineTransport fails every request, so tests never reach Discord.
type offlineTransport struct{}

func (offlineTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("offline")
}

// forbiddenTransport responds to every request with a Discord "Missing
// Permissions" error, and counts the requests.
type forbiddenTransport struct {
	requests atomic.Int64
}

func (f *forbiddenTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	f.requests.Add(1)
	return &http.Response{
		StatusCode: http.StatusForbidden,
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"code": 50013, "message": "Missing Permissions"}`)),
		Request:    r,
	}, nil
}

// newTestWorker returns a worker whose Discord REST client fails every
// request, and which is marked as running (as if [Worker.Run] was called,
// without connecting to the gateway) until ctx is cancelled. Gateway events
// can be simulated by calling the shard's handlers.
func newTestWorker(ctx context.Context, t *testing.T, db *ent.Client, config models.ConfigWorker) *Worker {
	t.Helper()
	return newTestWorkerTransport(ctx, t, db, config, offlineTransport{})
}

// newTestWorkerTransport is [newTestWorker], with a custom transport for the
// Discord REST client.
func newTestWorkerTransport(ctx context.Context, t *testing.T, db *ent.Client, config models.ConfigWorker, rt http.RoundTripper) *Worker {
	t.Helper()

	config.BotToken = "dummy"
	if config.EventFlushInterval == 0 {
		config.EventFlushInterval = 50 * time.Millisecond
	}

	w, err := New(ctx, slog.New(slog.DiscardHandler), db, config)
	if err != nil {
		t.Fatal(err)
	}
	w.client.Client.Client = httpdriver.WrapClient(http.Client{Transport: rt})
	w.client.Retries = 1

	w.mu.Lock()
	w.runCtx = ctx
	w.mu.Unlock()
	t.Cleanup(w.wg.Wait)
	return w
}

// fakeDiscord is a fake Discord API for a single guild, which applies changes
// to its own list of channels (so plans can be applied, and re-planned against
// the result), and records the calls made.
type fakeDiscord struct {
	guildID  discord.GuildID
	channels []discord.Channel
	nextID   discord.ChannelID

	creates  []api.CreateChannelData
	deletes  []discord.ChannelID
	modifies map[discord.ChannelID]api.ModifyChannelData
	moves    []api.MoveChannelsData

	// err, if set, returns an error for the given call (and channel ID, where
	// applicable, or 0 for any).
	err func(call string, id discord.ChannelID) error
}

var _ channelAPI = (*fakeDiscord)(nil)

func newFakeDiscord(channels ...discord.Channel) *fakeDiscord {
	return &fakeDiscord{
		guildID:  testGuildID,
		channels: slices.Clone(channels),
		nextID:   9000,
		modifies: make(map[discord.ChannelID]api.ModifyChannelData),
	}
}

func (f *fakeDiscord) fail(call string, id discord.ChannelID) error {
	if f.err == nil {
		return nil
	}
	return f.err(call, id)
}

func (f *fakeDiscord) index(id discord.ChannelID) int {
	return slices.IndexFunc(f.channels, func(ch discord.Channel) bool { return ch.ID == id })
}

func (f *fakeDiscord) CreateChannel(guildID discord.GuildID, data api.CreateChannelData) (*discord.Channel, error) {
	if err := f.fail("create", 0); err != nil {
		return nil, err
	}
	f.creates = append(f.creates, data)
	f.nextID++

	ch := discord.Channel{
		ID:             f.nextID,
		GuildID:        guildID,
		Type:           data.Type,
		Name:           data.Name,
		ParentID:       data.CategoryID,
		VoiceBitrate:   data.VoiceBitrate,
		VoiceUserLimit: data.VoiceUserLimit,
		Overwrites:     data.Overwrites,
		RTCRegionID:    data.RTCRegionID,
	}
	if data.Position != nil {
		ch.Position = *data.Position
	}
	f.channels = append(f.channels, ch)
	return &ch, nil
}

func (f *fakeDiscord) DeleteChannel(id discord.ChannelID, _ api.AuditLogReason) error {
	if err := f.fail("delete", id); err != nil {
		return err
	}
	f.deletes = append(f.deletes, id)
	if i := f.index(id); i >= 0 {
		f.channels = slices.Delete(f.channels, i, i+1)
	}
	return nil
}

func (f *fakeDiscord) ModifyChannel(id discord.ChannelID, data api.ModifyChannelData) error {
	if err := f.fail("modify", id); err != nil {
		return err
	}
	f.modifies[id] = data
	if i := f.index(id); i >= 0 {
		if data.VoiceBitrate != nil {
			f.channels[i].VoiceBitrate = data.VoiceBitrate.Val
		}
		if data.VoiceUserLimit != nil {
			f.channels[i].VoiceUserLimit = data.VoiceUserLimit.Val
		}
		if data.Overwrites != nil {
			f.channels[i].Overwrites = *data.Overwrites
		}
	}
	return nil
}

func (f *fakeDiscord) MoveChannels(_ discord.GuildID, data api.MoveChannelsData) error {
	if err := f.fail("move", 0); err != nil {
		return err
	}
	f.moves = append(f.moves, data)
	for _, m := range data.Channels {
		if i := f.index(m.ID); i >= 0 && m.Position != nil {
			f.channels[i].Position = *m.Position
		}
	}
	return nil
}

// fakeStore is a fake channel cache.
type fakeStore struct {
	set     []discord.ChannelID
	removed []discord.ChannelID

	// voiceStates are the (current) voice states of the guild.
	voiceStates []discord.VoiceState
}

var _ channelStore = (*fakeStore)(nil)

func (s *fakeStore) ChannelSet(ch *discord.Channel, _ bool) error {
	s.set = append(s.set, ch.ID)
	return nil
}

func (s *fakeStore) ChannelRemove(ch *discord.Channel) error {
	s.removed = append(s.removed, ch.ID)
	return nil
}

func (s *fakeStore) VoiceStates(discord.GuildID) ([]discord.VoiceState, error) {
	return s.voiceStates, nil
}
