// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package voicestate

import (
	"sync"

	"github.com/andersfylling/disgord"
)

// NewTracker returns a new voice state tracker.
func NewTracker() *Tracker {
	return &Tracker{
		db: make(map[disgord.Snowflake]map[disgord.Snowflake]*disgord.VoiceState),
	}
}

type Tracker struct {
	mu sync.RWMutex
	db map[disgord.Snowflake]map[disgord.Snowflake]*disgord.VoiceState
}

// Register should be used before the connection is established. Registers
// the necessary handlers for tracking voice state changes.
func (t *Tracker) Register(gw disgord.GatewayQueryBuilder) {
	gw.GuildCreate(func(s disgord.Session, h *disgord.GuildCreate) {
		t.process(h.Guild.ID, h.Guild.VoiceStates...)
	})

	gw.VoiceStateUpdate(func(s disgord.Session, h *disgord.VoiceStateUpdate) {
		t.process(h.GuildID, h.VoiceState)
	})

	gw.GuildDelete(func(s disgord.Session, h *disgord.GuildDelete) {
		if h.UserWasRemoved() {
			t.removeGuild(h.UnavailableGuild.ID)
		}
	})
}

func (t *Tracker) process(guildID disgord.Snowflake, states ...*disgord.VoiceState) {
	if states == nil {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	for _, state := range states {
		if _, ok := t.db[guildID]; !ok {
			t.db[guildID] = make(map[disgord.Snowflake]*disgord.VoiceState)
		}

		// https://discord.com/developers/docs/topics/gateway#update-voice-state
		//   channel_id: id of the voice channel client wants to join (null if disconnecting)
		if state.ChannelID.IsZero() {
			delete(t.db[guildID], state.UserID)
			continue
		}

		t.db[guildID][state.UserID] = state
	}
}

func (t *Tracker) removeGuild(guildID disgord.Snowflake) {
	if guildID.IsZero() {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.db, guildID)
}

// States returns the full list of known voice states for a given guild.
func (t *Tracker) States(guildID disgord.Snowflake) (states []*disgord.VoiceState) {
	if guildID.IsZero() {
		return states
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	if _, ok := t.db[guildID]; !ok {
		return states
	}

	states = make([]*disgord.VoiceState, 0, len(t.db[guildID]))
	for _, state := range t.db[guildID] {
		states = append(states, state)
	}

	return states
}

// UserCount returns a map where the keys are channels that have active voice states,
// and the value is the number of users in that voice channel, for a given guild.
func (t *Tracker) UserCount(guildID disgord.Snowflake) map[disgord.Snowflake]int {
	voiceCount := map[disgord.Snowflake]int{}
	for _, state := range t.States(guildID) {
		voiceCount[state.ChannelID]++
	}

	return voiceCount
}
