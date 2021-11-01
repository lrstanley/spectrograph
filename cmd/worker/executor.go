// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/apex/log"
)

type updateEvent struct {
	sess  disgord.Session
	event interface{}
	guild *disgord.Guild
}

func (b *discordBot) routeEvent(sess disgord.Session, event interface{}, guildID disgord.Snowflake) {
	update := &updateEvent{sess: sess, event: event}

	var err error
	update.guild, err = update.sess.Guild(guildID).Get()
	if err != nil {
		logGuild(logger, guildID).WithError(err).Error("unable to fetch guild for update event")
		return
	}

	b.updateMu.Lock()
	ch, ok := b.updates[guildID]
	if !ok {
		logGuild(logger, update.guild).Debug("dropping update event due to untracked guild")
		return
	}

	select {
	case ch <- update:
		b.updateMu.Unlock()
	default:
		b.updateMu.Unlock()
		// TODO: better way to handle us triggering subsequent events, as we
		// go through and make changes? i.e. us moving channels around will
		// cause events. Could we check the author?
		//
		// Alternatively, could we keep the last message, but do a wait and retry?
		logGuild(logger, update.guild).Debug("dropping event, as already processing event")
	}
}

func (b *discordBot) eventWatcher(sess disgord.Session, guildID disgord.Snowflake, events <-chan *updateEvent) {
	logGuild(logger, guildID).Debug("starting worker")
	defer logGuild(logger, guildID).Debug("closing worker")

	var event *updateEvent
	var ok bool

	// Debounce by this timeframe. I.e. if we receive multiple events within
	// the below timeframe, wait until we receive no new events within the
	// timeframe, then send out the last event we received.
	timer := time.NewTimer(5000 * time.Millisecond)

	for {
		select {
		case <-b.ctx.Done():
			return
		case <-timer.C:
			if event != nil {
				b.processUpdateWorker(sess, guildID, event) // TODO: if this panics, catch it.
				event = nil
			}
		case event, ok = <-events:
			if !ok {
				// Assume channel was closed because we were disconnected from
				// the guild.
				return
			}
			timer.Reset(500 * time.Millisecond)
		}
	}
}

func (b *discordBot) processUpdateWorker(sess disgord.Session, guildID disgord.Snowflake, event *updateEvent) {
	logGuild(logger, event.guild).Info("processing event")
	fmt.Printf("%#v\n", event)

	// Get server options.
	serverOptionsAdmin, err := svcServers.GetOptionsAdmin(b.ctx, event.guild.ID.String())
	if err != nil {
		logGuild(logger, event.guild).WithError(err).Error("unable to fetch server options (admin)")
		return
	}
	serverOptions, err := svcServers.GetOptions(b.ctx, event.guild.ID.String())
	if err != nil {
		logGuild(logger, event.guild).WithError(err).Error("unable to fetch server options")
		return
	}

	if !serverOptionsAdmin.Enabled || !serverOptions.Enabled {
		logGuild(logger, event.guild).Debug("dropping event, guild disabled via options")
		return
	}

	// TODO: make a cache for this (with key just being the regex string).
	// Probably TTL it or key off the guild ID so it can't be abused.
	rgx, err := regexp.Compile(serverOptions.RegexMatch)
	if err != nil {
		logGuild(logger, event.guild).WithError(err).WithField("regex", serverOptions.RegexMatch).Error("unable to parse regex") // TODO
		return
	}

	// Get number of users in each voice channel.
	voiceCount := b.voiceStates.UserCount(event.guild.ID)

	// state defines the state of parent channels and their "buckets" of managed
	// channels. See also:
	//   {
	//   	parent.Snowflake: {
	//   		channel1.Name: []<channel>,
	//   		channel2.Name: []<channel>,
	//   	}
	//   }
	state := map[disgord.Snowflake]map[string][]*disgord.Channel{}
	for _, channel := range event.guild.Channels {
		if channel.Type != disgord.ChannelTypeGuildVoice {
			continue
		}

		if ok := rgx.MatchString(channel.Name); !ok {
			continue
		}

		if _, ok := state[channel.ParentID]; !ok {
			state[channel.ParentID] = map[string][]*disgord.Channel{}
		}

		if _, ok := state[channel.ParentID][channel.Name]; !ok {
			state[channel.ParentID][channel.Name] = []*disgord.Channel{channel}
		} else {
			state[channel.ParentID][channel.Name] = append(state[channel.ParentID][channel.Name], channel)
		}
	}

	// TODO: count how many channel groups we have, as well as how many channels
	// in each group. Check this against our configured limits.

	// TODO: should we make sure we have permissions?
	// TODO: support empty channel being at the top, vs the bottom.
	// TODO: support multiple empty channels? could help in the event of bot
	// issues if there is a "buffer" of empty channels.

	var emptyChannel *disgord.Channel
	var lastOccupiedChannel *disgord.Channel

	toDelete := []*disgord.Channel{}

	for parent := range state {
		for group := range state[parent] {
			emptyChannel = nil
			lastOccupiedChannel = nil

			// Find which channels are empty, and which have users. If there
			// are more than one channels that are empty, mark all subsequent
			// empty channels for deletion.
			for _, channel := range state[parent][group] {
				if voiceCount[channel.ID] == 0 {
					if emptyChannel != nil {
						toDelete = append(toDelete, channel)
						continue
					}

					emptyChannel = channel
				} else {
					lastOccupiedChannel = channel
				}
			}

			// Move empty channel to position after the last occupied channel.
			if emptyChannel != nil && lastOccupiedChannel != nil && lastOccupiedChannel.ID != emptyChannel.ID {
				err = sess.Guild(event.guild.ID).UpdateChannelPositions([]disgord.UpdateGuildChannelPositionsParams{
					disgord.UpdateGuildChannelPositionsParams{ID: lastOccupiedChannel.ID, Position: lastOccupiedChannel.Position},
					disgord.UpdateGuildChannelPositionsParams{ID: emptyChannel.ID, Position: lastOccupiedChannel.Position + 1},
				})
				if err != nil {
					logGuild(logger, event.guild).WithError(err).WithFields(log.Fields{
						"last_occupied_id": lastOccupiedChannel.ID,
						"empty_channel_id": emptyChannel.ID,
					}).Error("unable to reorder empty channel to bottom")
				}
			}

			// If no empty channel, make one, duplicating the config from the first
			// channel in the bucket.
			if emptyChannel == nil {
				channel, err := sess.Guild(event.guild.ID).CreateChannel(state[parent][group][0].Name, &disgord.CreateGuildChannelParams{
					Name:                 state[parent][group][0].Name,
					Type:                 state[parent][group][0].Type,
					Bitrate:              state[parent][group][0].Bitrate,
					UserLimit:            state[parent][group][0].UserLimit,
					RateLimitPerUser:     state[parent][group][0].RateLimitPerUser,
					PermissionOverwrites: state[parent][group][0].PermissionOverwrites,
					ParentID:             state[parent][group][0].ParentID,
					NSFW:                 state[parent][group][0].NSFW,
					Position:             state[parent][group][0].Position + 1,
				})
				if err != nil {
					logGuild(logger, event.guild).WithError(err).WithField("source_channel_id", state[parent][group][0]).Error("unable to create new channel from master channel")
				} else {
					emptyChannel = channel

					// add new channel to state.
					state[parent][group] = append(state[parent][group], channel)
				}
			}
		}
	}

	for _, channel := range toDelete {
		if _, err := sess.Channel(channel.ID).Delete(); err != nil {
			// This can sometimes cause "unknown channel" if executed too fast
			// between other API calls.
			if restErr, ok := err.(*disgord.ErrRest); ok {
				if restErr.Code == 10003 {
					// https://discord.com/developers/docs/topics/opcodes-and-status-codes#json-json-error-codes
					continue
				}
			}

			logGuild(logger, event.guild).WithError(err).WithField("channel_id", channel.ID).Error("unable to remove unneeded empty channel")
		}

		// Remove the deleted channel from state.
		group := state[channel.ParentID][channel.Name]
		for i := 0; i < len(group); i++ {
			if group[i].Compare(channel) {
				copy(group[i:], group[i+1:])
				group[len(group)-1] = nil                                    // remove last element to prevent memory leaking.
				state[channel.ParentID][channel.Name] = group[:len(group)-1] // truncate slice.

				break
			}
		}
	}

	// Loop through all of the channels and make sure their config matches that of the
	// "primary" channel in the list. I.e. change the primary, and the rest should
	// change.
	for parent := range state {
		// parentChannel, err := sess.Channel(parent).Get()
		// if err != nil {
		// 	pretty.Println(err)
		// } else {
		// 	pretty.Println(parentChannel)
		// }
		for group := range state[parent] {
			// Check if it's just one channel.
			if len(state[parent][group]) < 2 {
				continue
			}

			primary := state[parent][group][0]
			var needsUpdate bool

			for _, channel := range state[parent][group] {
				if channel.UserLimit != primary.UserLimit ||
					channel.Bitrate != primary.Bitrate ||
					len(channel.PermissionOverwrites) != len(primary.PermissionOverwrites) {
					needsUpdate = true
				}

				if !needsUpdate {
					for i := 0; i < len(channel.PermissionOverwrites); i++ {
						if channel.PermissionOverwrites[i].Type != primary.PermissionOverwrites[i].Type ||
							channel.PermissionOverwrites[i].ID != primary.PermissionOverwrites[i].ID ||
							channel.PermissionOverwrites[i].Allow != primary.PermissionOverwrites[i].Allow ||
							channel.PermissionOverwrites[i].Deny != primary.PermissionOverwrites[i].Deny {
							needsUpdate = true
							break
						}
					}
				}

				if !needsUpdate {
					continue
				}

				_, err := sess.Channel(channel.ID).UpdateBuilder().
					SetPosition(channel.Position).
					SetUserLimit(primary.UserLimit).
					SetBitrate(primary.Bitrate).
					SetPermissionOverwrites(primary.PermissionOverwrites).Execute()
				if err != nil {
					// TODO: this should be propagated up to the user somehow. events?
					// TODO: should we change the permissions ourselves?
					logGuild(logger, event.guild).WithError(err).WithFields(log.Fields{
						"channel_id": channel.ID,
						"primary_id": primary.ID,
					}).Error("unable to update children channel details based off primary channel")
				}
			}
		}
	}
}

// TODO: function to re-order and/or add permissions specifically for the bot user, into
// the channel permission overrides, ONLY if there is a permission that disallows
// being able to read/update, etc??
func (b *discordBot) changeChannelPermissions(sess *disgord.Session, channel *disgord.Channel) error {
	return nil
}
