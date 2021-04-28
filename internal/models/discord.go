// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.
//
// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by the following license:
//   https://github.com/bwmarrin/discordgo/blob/master/LICENSE
//
// Pulled from:
// https://github.com/bwmarrin/discordgo/blob/f7db9886fc14d6af29a19d38b647b489af8cdcd4/structs.go#L1232-L1312

package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// DiscordPermissions is used to convert Discord permission bits to Go bitwise-able
// values for comparison and validation.
// Discord docs: https://discord.com/developers/docs/topics/permissions#permissions-bitwise-permission-flags
type DiscordPermissions uint64

// Do type validation against the DiscordPermissions struct to ensure it matches
// the Marshaler/Unmarshaler interfaces.
var _ json.Marshaler = (*DiscordPermissions)(nil)
var _ json.Unmarshaler = (*DiscordPermissions)(nil)

func (b *DiscordPermissions) MarshalJSON() ([]byte, error) {
	str := strconv.FormatUint(uint64(*b), 10)
	return []byte(strconv.Quote(str)), nil
}

func (b *DiscordPermissions) UnmarshalJSON(bytes []byte) error {
	sb := string(bytes)
	str, err := strconv.Unquote(sb)
	if err != nil {
		return fmt.Errorf("unable to convert permission %v into type: %w", sb, err)
	}

	v, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return fmt.Errorf("unable to parse permission string to uint64 failed: %w", err)
	}

	*b = DiscordPermissions(v)
	return nil
}

// Contains is used to check if the permission integer contains the bits specified.
func (b DiscordPermissions) Contains(Bits DiscordPermissions) bool {
	return (b & Bits) == Bits
}

// Constants for the different bit offsets of text channel permissions.
const (
	// Deprecated: DiscordPermReadMessages has been replaced with DiscordPermViewChannel for text and voice channels.
	DiscordPermReadMessages       = 0x0000000000000400
	DiscordPermSendMessages       = 0x0000000000000800
	DiscordPermSendTTSMessages    = 0x0000000000001000
	DiscordPermManageMessages     = 0x0000000000002000
	DiscordPermEmbedLinks         = 0x0000000000004000
	DiscordPermAttachFiles        = 0x0000000000008000
	DiscordPermReadMessageHistory = 0x0000000000010000
	DiscordPermMentionEveryone    = 0x0000000000020000
	DiscordPermUseExternalEmojis  = 0x0000000000040000
	DiscordPermUseSlashCommands   = 0x0000000080000000
)

// Constants for the different bit offsets of voice permissions.
const (
	DiscordPermVoicePrioritySpeaker = 0x0000000000000100
	DiscordPermVoiceStreamVideo     = 0x0000000000000200
	DiscordPermVoiceConnect         = 0x0000000000100000
	DiscordPermVoiceSpeak           = 0x0000000000200000
	DiscordPermVoiceMuteMembers     = 0x0000000000400000
	DiscordPermVoiceDeafenMembers   = 0x0000000000800000
	DiscordPermVoiceMoveMembers     = 0x0000000001000000
	DiscordPermVoiceUseVAD          = 0x0000000002000000
	DiscordPermVoiceRequestToSpeak  = 0x0000000100000000
)

// Constants for general management.
const (
	DiscordPermChangeNickname  = 0x0000000004000000
	DiscordPermManageNicknames = 0x0000000008000000
	DiscordPermManageRoles     = 0x0000000010000000
	DiscordPermManageWebhooks  = 0x0000000020000000
	DiscordPermManageEmojis    = 0x0000000040000000
)

// Constants for the different bit offsets of general permissions.
const (
	DiscordPermCreateInstantInvite = 0x0000000000000001
	DiscordPermKickMembers         = 0x0000000000000002
	DiscordPermBanMembers          = 0x0000000000000004
	DiscordPermAdministrator       = 0x0000000000000008
	DiscordPermManageChannels      = 0x0000000000000010
	DiscordPermManageServer        = 0x0000000000000020
	DiscordPermAddReactions        = 0x0000000000000040
	DiscordPermViewAuditLogs       = 0x0000000000000080
	DiscordPermViewChannel         = 0x0000000000000400
	DiscordPermViewGuildInsights   = 0x0000000000080000

	DiscordPermAllText = DiscordPermViewChannel |
		DiscordPermSendMessages |
		DiscordPermSendTTSMessages |
		DiscordPermManageMessages |
		DiscordPermEmbedLinks |
		DiscordPermAttachFiles |
		DiscordPermReadMessageHistory |
		DiscordPermMentionEveryone
	DiscordPermAllVoice = DiscordPermViewChannel |
		DiscordPermVoiceConnect |
		DiscordPermVoiceSpeak |
		DiscordPermVoiceMuteMembers |
		DiscordPermVoiceDeafenMembers |
		DiscordPermVoiceMoveMembers |
		DiscordPermVoiceUseVAD |
		DiscordPermVoicePrioritySpeaker
	DiscordPermAllChannel = DiscordPermAllText |
		DiscordPermAllVoice |
		DiscordPermCreateInstantInvite |
		DiscordPermManageRoles |
		DiscordPermManageChannels |
		DiscordPermAddReactions |
		DiscordPermViewAuditLogs
	DiscordPermAll = DiscordPermAllChannel |
		DiscordPermKickMembers |
		DiscordPermBanMembers |
		DiscordPermManageServer |
		DiscordPermAdministrator |
		DiscordPermManageWebhooks |
		DiscordPermManageEmojis
)
