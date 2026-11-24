// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
)

func discordErr(status int, code httputil.ErrorCode) error {
	// arikawa returns *httputil.HTTPError.
	var err error = &httputil.HTTPError{Status: status, Code: code, Message: "test"}
	return fmt.Errorf("wrapped: %w", err)
}

func TestApplyDeleteUnknownChannel(t *testing.T) {
	t.Parallel()

	fake := newFakeDiscord(voice(1, "Gaming +", 1), voice(2, "Gaming +", 2), voice(3, "Gaming +", 3))
	fake.err = func(call string, id discord.ChannelID) error {
		switch {
		case call == "delete" && id == 2:
			return discordErr(http.StatusNotFound, errCodeUnknownChannel)
		case call == "delete" && id == 3:
			return discordErr(http.StatusForbidden, 50013)
		}
		return nil
	}

	cache := &fakeStore{}
	problems := apply(fake, cache, testGuildID, &plan{deletes: []discord.Channel{voice(2, "Gaming +", 2), voice(3, "Gaming +", 3)}})

	// Unknown channels are already gone, so they are removed from the cache
	// (and not reported), other errors are reported, and the channel kept.
	if !slices.Equal(cache.removed, []discord.ChannelID{2}) {
		t.Errorf("removed from cache = %v, want [2]", cache.removed)
	}
	if len(problems) != 1 || problems[0].fields[fieldChannelID] != "3" {
		t.Errorf("problems = %v, want 1 problem for channel 3", problems)
	}
}

// TestApplyDeleteSkipsOccupied checks that a channel planned for deletion isn't
// deleted if a user joined it since the plan was made.
func TestApplyDeleteSkipsOccupied(t *testing.T) {
	t.Parallel()

	fake := newFakeDiscord(voice(1, "Gaming +", 1), voice(2, "Gaming +", 2), voice(3, "Gaming +", 3))
	cache := &fakeStore{voiceStates: []discord.VoiceState{{GuildID: testGuildID, ChannelID: 2, UserID: testOtherUser}}}

	problems := apply(fake, cache, testGuildID, &plan{deletes: []discord.Channel{voice(2, "Gaming +", 2), voice(3, "Gaming +", 3)}})
	if len(problems) > 0 {
		t.Errorf("problems = %v, want none", problems)
	}
	if !slices.Equal(fake.deletes, []discord.ChannelID{3}) {
		t.Errorf("deleted = %v, want [3]", fake.deletes)
	}
}

func TestApplyCreate(t *testing.T) {
	t.Parallel()

	src := inCategory(voice(1, "Gaming +", 1), testCategory)
	src.Overwrites = []discord.Overwrite{
		{ID: discord.Snowflake(everyoneRole), Type: discord.OverwriteRole, Deny: discord.PermissionConnect},
	}
	src.RTCRegionID = "us-east"

	fake := newFakeDiscord(src)
	cache := &fakeStore{}

	problems := apply(fake, cache, testGuildID, &plan{creates: []createAction{{source: src, position: 2, overwrites: src.Overwrites}}})
	if len(problems) > 0 {
		t.Fatalf("problems = %v", problems)
	}

	if len(fake.creates) != 1 {
		t.Fatalf("creates = %d, want 1", len(fake.creates))
	}

	got := fake.creates[0]
	if got.Name != src.Name || got.Type != discord.GuildVoice || got.CategoryID != testCategory ||
		got.VoiceBitrate != src.VoiceBitrate || got.VoiceUserLimit != src.VoiceUserLimit ||
		got.RTCRegionID != "us-east" || got.Position == nil || *got.Position != 2 ||
		!overwritesEqual(got.Overwrites, src.Overwrites) {
		t.Errorf("unexpected create data: %+v", got)
	}

	// Created channels are cached immediately, so the next pass sees them.
	if !slices.Equal(cache.set, []discord.ChannelID{fake.nextID}) {
		t.Errorf("cached = %v, want [%d]", cache.set, fake.nextID)
	}
}

func TestApplyCreateError(t *testing.T) {
	t.Parallel()

	fake := newFakeDiscord()
	fake.err = func(string, discord.ChannelID) error { return errors.New("boom") }
	cache := &fakeStore{}

	problems := apply(fake, cache, testGuildID, &plan{creates: []createAction{{source: voice(1, "Gaming +", 1), position: 2}}})
	if len(problems) != 1 || problems[0].fields["error"] != "boom" {
		t.Errorf("problems = %v, want 1 problem", problems)
	}
	if len(cache.set) > 0 {
		t.Errorf("cached = %v, want none", cache.set)
	}
}

func TestApplyMovesBatched(t *testing.T) {
	t.Parallel()

	fake := newFakeDiscord(voice(1, "A +", 1), voice(2, "B +", 2))
	problems := apply(fake, &fakeStore{}, testGuildID, &plan{moves: []moveAction{
		{channel: voice(1, "A +", 1), position: 5},
		{channel: voice(2, "B +", 2), position: 7},
	}})
	if len(problems) > 0 {
		t.Fatalf("problems = %v", problems)
	}

	if len(fake.moves) != 1 || len(fake.moves[0].Channels) != 2 {
		t.Fatalf("moves = %+v, want a single request with 2 channels", fake.moves)
	}
	for _, m := range fake.moves[0].Channels {
		if m.Position == nil {
			t.Errorf("move %d has no position", m.ID)
		}
	}
}

func TestApplySync(t *testing.T) {
	t.Parallel()

	primary := voice(1, "Gaming +", 1)
	primary.VoiceUserLimit = 3
	primary.Overwrites = []discord.Overwrite{
		{ID: discord.Snowflake(everyoneRole), Type: discord.OverwriteRole, Deny: discord.PermissionConnect},
		{ID: discord.Snowflake(testRoleID), Type: discord.OverwriteRole, Allow: discord.PermissionConnect},
	}

	clone := voice(2, "Gaming +", 2)
	clone.VoiceUserLimit = 0 // Unlimited, must still be sent.

	settingsOnly := voice(3, "Gaming +", 3)

	fake := newFakeDiscord(primary, clone, settingsOnly)
	problems := apply(fake, &fakeStore{}, testGuildID, &plan{syncs: []syncAction{
		{channel: clone, source: primary, overwrites: true, permissionOverwrites: primary.Overwrites},
		{channel: settingsOnly, source: primary},
	}})
	if len(problems) > 0 {
		t.Fatalf("problems = %v", problems)
	}

	got := fake.modifies[2]
	if got.VoiceUserLimit == nil || got.VoiceUserLimit.Val != 3 {
		t.Errorf("user limit = %+v, want 3", got.VoiceUserLimit)
	}
	if got.VoiceBitrate == nil || got.VoiceBitrate.Val != primary.VoiceBitrate {
		t.Errorf("bitrate = %+v, want %d", got.VoiceBitrate, primary.VoiceBitrate)
	}
	// Full overwrites (IDs, types, allow and deny) must be sent.
	if got.Overwrites == nil || !slices.Equal(*got.Overwrites, primary.Overwrites) {
		t.Errorf("overwrites = %+v, want %+v", got.Overwrites, primary.Overwrites)
	}

	if fake.modifies[3].Overwrites != nil {
		t.Errorf("overwrites sent for settings-only sync: %+v", *fake.modifies[3].Overwrites)
	}
}

func TestApplySyncClearsOverwrites(t *testing.T) {
	t.Parallel()

	primary := voice(1, "Gaming +", 1)
	clone := voice(2, "Gaming +", 2)
	clone.Overwrites = []discord.Overwrite{{ID: discord.Snowflake(everyoneRole), Type: discord.OverwriteRole, Deny: discord.PermissionConnect}}

	fake := newFakeDiscord(primary, clone)
	apply(fake, &fakeStore{}, testGuildID, &plan{syncs: []syncAction{{channel: clone, source: primary, overwrites: true}}})

	// An empty (non-nil) list must be sent, as nil would omit the field.
	if ow := fake.modifies[2].Overwrites; ow == nil || *ow == nil || len(*ow) != 0 {
		t.Errorf("overwrites = %v, want empty list", ow)
	}
}

func TestErrorHelpers(t *testing.T) {
	t.Parallel()

	if !isDiscordError(discordErr(http.StatusNotFound, errCodeUnknownChannel), errCodeUnknownChannel) {
		t.Error("isDiscordError() = false, want true")
	}
	if isDiscordError(errors.New("x"), errCodeUnknownChannel) {
		t.Error("isDiscordError() = true for non-discord error")
	}
	if !isUnauthorized(discordErr(http.StatusUnauthorized, 0)) {
		t.Error("isUnauthorized() = false, want true")
	}
	if isUnauthorized(discordErr(http.StatusForbidden, 50013)) {
		t.Error("isUnauthorized() = true for 403")
	}
}
