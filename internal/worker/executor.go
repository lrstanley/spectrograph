// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"errors"
	"log/slog"
	"net/http"
	"slices"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

// auditLogReason is shown in the guild audit log for changes made by the bot.
const auditLogReason api.AuditLogReason = "spectrograph: managing voice channel clones"

// Discord JSON error codes, see:
// https://discord.com/developers/docs/topics/opcodes-and-status-codes#json-json-error-codes
const errCodeUnknownChannel httputil.ErrorCode = 10003

// channelAPI is the subset of the Discord REST API used to apply plans.
// Implemented by [api.Client].
type channelAPI interface {
	CreateChannel(guildID discord.GuildID, data api.CreateChannelData) (*discord.Channel, error)
	DeleteChannel(channelID discord.ChannelID, reason api.AuditLogReason) error
	ModifyChannel(channelID discord.ChannelID, data api.ModifyChannelData) error
	MoveChannels(guildID discord.GuildID, data api.MoveChannelsData) error
}

// channelStore is the subset of the state cache updated with the results of
// applied changes (and checked for users who joined a channel since the plan
// was made). Implemented by [store.Cabinet].
type channelStore interface {
	ChannelSet(ch *discord.Channel, update bool) error
	ChannelRemove(ch *discord.Channel) error
	VoiceStates(guildID discord.GuildID) ([]discord.VoiceState, error)
}

// isOccupied returns true if the cache has a user connected to the channel.
// The cache is kept up to date by the gateway while a plan is applied.
func isOccupied(cache channelStore, guildID discord.GuildID, id discord.ChannelID) bool {
	states, _ := cache.VoiceStates(guildID)
	return slices.ContainsFunc(states, func(vs discord.VoiceState) bool { return vs.ChannelID == id })
}

// isDiscordError returns true if err is a Discord API error with the given
// JSON error code.
func isDiscordError(err error, code httputil.ErrorCode) bool {
	var herr *httputil.HTTPError
	return errors.As(err, &herr) && herr.Code == code
}

// isUnauthorized returns true if err is a Discord API 401 error (i.e. an
// invalid token).
func isUnauthorized(err error) bool {
	var herr *httputil.HTTPError
	return errors.As(err, &herr) && herr.Status == http.StatusUnauthorized
}

// isForbidden returns true if err is a Discord API 401 or 403 error, which
// count towards Discord's invalid request limit.
func isForbidden(err error) bool {
	var herr *httputil.HTTPError
	return errors.As(err, &herr) && (herr.Status == http.StatusUnauthorized || herr.Status == http.StatusForbidden)
}

// isTransient returns true if err isn't a Discord API client error (4xx, e.g.
// missing permissions or an invalid request, which would fail again, and count
// towards Discord's invalid request limit), so retrying may succeed.
func isTransient(err error) bool {
	if herr, ok := errors.AsType[*httputil.HTTPError](err); ok {
		return herr.Status >= http.StatusInternalServerError || herr.Status == http.StatusTooManyRequests
	}
	return true
}

// apply applies the plan to the guild, and returns any problems encountered.
// Successful creates and deletes are reflected in the cache immediately, so a
// subsequent pass doesn't act on stale data before the corresponding gateway
// events arrive (e.g. creating a second clone).
func apply(client channelAPI, cache channelStore, guildID discord.GuildID, p *plan) (problems []problem) {
	fail := func(msg string, err error, fields map[string]any) {
		fields["error"] = err.Error()
		problems = append(problems, problem{
			level:     slog.LevelError,
			msg:       msg,
			fields:    fields,
			retry:     isTransient(err),
			forbidden: isForbidden(err),
		})
	}

	for _, ch := range p.deletes {
		// A user may have joined the channel since the plan was made (or while
		// previous changes were applied). Deleting it would disconnect them,
		// and the voice state update triggers another pass anyway.
		if isOccupied(cache, guildID, ch.ID) {
			continue
		}

		err := client.DeleteChannel(ch.ID, auditLogReason)
		if err != nil && !isDiscordError(err, errCodeUnknownChannel) {
			fail("unable to remove unneeded empty channel", err, map[string]any{fieldChannelID: ch.ID.String()})
			continue
		}
		// Also remove channels Discord reports as unknown (already deleted), so
		// they don't linger in the cache.
		_ = cache.ChannelRemove(&ch)
	}

	for _, c := range p.creates {
		src := c.source
		ch, err := client.CreateChannel(guildID, api.CreateChannelData{
			Name:           src.Name,
			Type:           src.Type,
			VoiceBitrate:   src.VoiceBitrate,
			VoiceUserLimit: src.VoiceUserLimit,
			UserRateLimit:  src.UserRateLimit,
			Position:       option.NewInt(c.position),
			Overwrites:     c.overwrites,
			CategoryID:     src.ParentID,
			NSFW:           src.NSFW,
			RTCRegionID:    src.RTCRegionID,
			AuditLogReason: auditLogReason,
		})
		if err != nil {
			fail("unable to clone channel", err, map[string]any{fieldChannelID: src.ID.String()})
			continue
		}
		if !ch.GuildID.IsValid() {
			ch.GuildID = guildID
		}
		_ = cache.ChannelSet(ch, false)
	}

	if len(p.moves) > 0 {
		data := api.MoveChannelsData{AuditLogReason: auditLogReason}
		for _, m := range p.moves {
			data.Channels = append(data.Channels, api.MoveChannelData{ID: m.channel.ID, Position: option.NewInt(m.position)})
		}

		if err := client.MoveChannels(guildID, data); err != nil {
			ids := make([]string, 0, len(p.moves))
			for _, m := range p.moves {
				ids = append(ids, m.channel.ID.String())
			}
			fail("unable to move empty channels after occupied channels", err, map[string]any{"channel_ids": ids})
		}
	}

	for i := range p.syncs {
		s := &p.syncs[i]
		data := api.ModifyChannelData{
			VoiceBitrate:   option.NewNullableUint(s.source.VoiceBitrate),
			VoiceUserLimit: option.NewNullableUint(s.source.VoiceUserLimit),
			AuditLogReason: auditLogReason,
		}
		if s.overwrites {
			overwrites := s.permissionOverwrites
			if overwrites == nil {
				overwrites = []discord.Overwrite{}
			}
			data.Overwrites = &overwrites
		}

		if err := client.ModifyChannel(s.channel.ID, data); err != nil {
			fail("unable to sync clone with original channel", err, map[string]any{
				fieldChannelID: s.channel.ID.String(),
				fieldPrimaryID: s.source.ID.String(),
			})
		}
	}

	return problems
}
