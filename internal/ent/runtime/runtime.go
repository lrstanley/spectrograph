// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.
//
// DO NOT EDIT, CODE GENERATED BY entc.

package runtime

import (
	"context"
	"time"

	"github.com/lrstanley/spectrograph/internal/database/schema"
	"github.com/lrstanley/spectrograph/internal/ent/guild"
	"github.com/lrstanley/spectrograph/internal/ent/guildadminconfig"
	"github.com/lrstanley/spectrograph/internal/ent/guildconfig"
	"github.com/lrstanley/spectrograph/internal/ent/guildevent"
	"github.com/lrstanley/spectrograph/internal/ent/user"

	"entgo.io/ent"
	"entgo.io/ent/privacy"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	guildMixin := schema.Guild{}.Mixin()
	guild.Policy = privacy.NewPolicies(schema.Guild{})
	guild.Hooks[0] = func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if err := guild.Policy.EvalMutation(ctx, m); err != nil {
				return nil, err
			}
			return next.Mutate(ctx, m)
		})
	}
	guildMixinFields0 := guildMixin[0].Fields()
	_ = guildMixinFields0
	guildFields := schema.Guild{}.Fields()
	_ = guildFields
	// guildDescCreateTime is the schema descriptor for create_time field.
	guildDescCreateTime := guildMixinFields0[0].Descriptor()
	// guild.DefaultCreateTime holds the default value on creation for the create_time field.
	guild.DefaultCreateTime = guildDescCreateTime.Default.(func() time.Time)
	// guildDescUpdateTime is the schema descriptor for update_time field.
	guildDescUpdateTime := guildMixinFields0[1].Descriptor()
	// guild.DefaultUpdateTime holds the default value on creation for the update_time field.
	guild.DefaultUpdateTime = guildDescUpdateTime.Default.(func() time.Time)
	// guild.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	guild.UpdateDefaultUpdateTime = guildDescUpdateTime.UpdateDefault.(func() time.Time)
	// guildDescName is the schema descriptor for name field.
	guildDescName := guildFields[1].Descriptor()
	// guild.NameValidator is a validator for the "name" field. It is called by the builders before save.
	guild.NameValidator = func() func(string) error {
		validators := guildDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// guildDescFeatures is the schema descriptor for features field.
	guildDescFeatures := guildFields[2].Descriptor()
	// guild.DefaultFeatures holds the default value on creation for the features field.
	guild.DefaultFeatures = guildDescFeatures.Default.([]string)
	// guildDescIconHash is the schema descriptor for icon_hash field.
	guildDescIconHash := guildFields[3].Descriptor()
	// guild.IconHashValidator is a validator for the "icon_hash" field. It is called by the builders before save.
	guild.IconHashValidator = guildDescIconHash.Validators[0].(func(string) error)
	// guildDescIconURL is the schema descriptor for icon_url field.
	guildDescIconURL := guildFields[4].Descriptor()
	// guild.IconURLValidator is a validator for the "icon_url" field. It is called by the builders before save.
	guild.IconURLValidator = guildDescIconURL.Validators[0].(func(string) error)
	// guildDescLarge is the schema descriptor for large field.
	guildDescLarge := guildFields[6].Descriptor()
	// guild.DefaultLarge holds the default value on creation for the large field.
	guild.DefaultLarge = guildDescLarge.Default.(bool)
	// guildDescMemberCount is the schema descriptor for member_count field.
	guildDescMemberCount := guildFields[7].Descriptor()
	// guild.DefaultMemberCount holds the default value on creation for the member_count field.
	guild.DefaultMemberCount = guildDescMemberCount.Default.(int)
	// guildDescPermissions is the schema descriptor for permissions field.
	guildDescPermissions := guildFields[9].Descriptor()
	// guild.DefaultPermissions holds the default value on creation for the permissions field.
	guild.DefaultPermissions = guildDescPermissions.Default.(uint64)
	// guildDescSystemChannelFlags is the schema descriptor for system_channel_flags field.
	guildDescSystemChannelFlags := guildFields[10].Descriptor()
	// guild.SystemChannelFlagsValidator is a validator for the "system_channel_flags" field. It is called by the builders before save.
	guild.SystemChannelFlagsValidator = guildDescSystemChannelFlags.Validators[0].(func(string) error)
	guildadminconfigMixin := schema.GuildAdminConfig{}.Mixin()
	guildadminconfig.Policy = privacy.NewPolicies(schema.GuildAdminConfig{})
	guildadminconfig.Hooks[0] = func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if err := guildadminconfig.Policy.EvalMutation(ctx, m); err != nil {
				return nil, err
			}
			return next.Mutate(ctx, m)
		})
	}
	guildadminconfigMixinFields0 := guildadminconfigMixin[0].Fields()
	_ = guildadminconfigMixinFields0
	guildadminconfigFields := schema.GuildAdminConfig{}.Fields()
	_ = guildadminconfigFields
	// guildadminconfigDescCreateTime is the schema descriptor for create_time field.
	guildadminconfigDescCreateTime := guildadminconfigMixinFields0[0].Descriptor()
	// guildadminconfig.DefaultCreateTime holds the default value on creation for the create_time field.
	guildadminconfig.DefaultCreateTime = guildadminconfigDescCreateTime.Default.(func() time.Time)
	// guildadminconfigDescUpdateTime is the schema descriptor for update_time field.
	guildadminconfigDescUpdateTime := guildadminconfigMixinFields0[1].Descriptor()
	// guildadminconfig.DefaultUpdateTime holds the default value on creation for the update_time field.
	guildadminconfig.DefaultUpdateTime = guildadminconfigDescUpdateTime.Default.(func() time.Time)
	// guildadminconfig.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	guildadminconfig.UpdateDefaultUpdateTime = guildadminconfigDescUpdateTime.UpdateDefault.(func() time.Time)
	// guildadminconfigDescEnabled is the schema descriptor for enabled field.
	guildadminconfigDescEnabled := guildadminconfigFields[0].Descriptor()
	// guildadminconfig.DefaultEnabled holds the default value on creation for the enabled field.
	guildadminconfig.DefaultEnabled = guildadminconfigDescEnabled.Default.(bool)
	// guildadminconfigDescDefaultMaxChannels is the schema descriptor for default_max_channels field.
	guildadminconfigDescDefaultMaxChannels := guildadminconfigFields[1].Descriptor()
	// guildadminconfig.DefaultDefaultMaxChannels holds the default value on creation for the default_max_channels field.
	guildadminconfig.DefaultDefaultMaxChannels = guildadminconfigDescDefaultMaxChannels.Default.(int)
	// guildadminconfig.DefaultMaxChannelsValidator is a validator for the "default_max_channels" field. It is called by the builders before save.
	guildadminconfig.DefaultMaxChannelsValidator = guildadminconfigDescDefaultMaxChannels.Validators[0].(func(int) error)
	// guildadminconfigDescDefaultMaxClones is the schema descriptor for default_max_clones field.
	guildadminconfigDescDefaultMaxClones := guildadminconfigFields[2].Descriptor()
	// guildadminconfig.DefaultDefaultMaxClones holds the default value on creation for the default_max_clones field.
	guildadminconfig.DefaultDefaultMaxClones = guildadminconfigDescDefaultMaxClones.Default.(int)
	// guildadminconfig.DefaultMaxClonesValidator is a validator for the "default_max_clones" field. It is called by the builders before save.
	guildadminconfig.DefaultMaxClonesValidator = guildadminconfigDescDefaultMaxClones.Validators[0].(func(int) error)
	// guildadminconfigDescComment is the schema descriptor for comment field.
	guildadminconfigDescComment := guildadminconfigFields[3].Descriptor()
	// guildadminconfig.DefaultComment holds the default value on creation for the comment field.
	guildadminconfig.DefaultComment = guildadminconfigDescComment.Default.(string)
	guildconfigMixin := schema.GuildConfig{}.Mixin()
	guildconfig.Policy = privacy.NewPolicies(schema.GuildConfig{})
	guildconfig.Hooks[0] = func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if err := guildconfig.Policy.EvalMutation(ctx, m); err != nil {
				return nil, err
			}
			return next.Mutate(ctx, m)
		})
	}
	guildconfigMixinFields0 := guildconfigMixin[0].Fields()
	_ = guildconfigMixinFields0
	guildconfigFields := schema.GuildConfig{}.Fields()
	_ = guildconfigFields
	// guildconfigDescCreateTime is the schema descriptor for create_time field.
	guildconfigDescCreateTime := guildconfigMixinFields0[0].Descriptor()
	// guildconfig.DefaultCreateTime holds the default value on creation for the create_time field.
	guildconfig.DefaultCreateTime = guildconfigDescCreateTime.Default.(func() time.Time)
	// guildconfigDescUpdateTime is the schema descriptor for update_time field.
	guildconfigDescUpdateTime := guildconfigMixinFields0[1].Descriptor()
	// guildconfig.DefaultUpdateTime holds the default value on creation for the update_time field.
	guildconfig.DefaultUpdateTime = guildconfigDescUpdateTime.Default.(func() time.Time)
	// guildconfig.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	guildconfig.UpdateDefaultUpdateTime = guildconfigDescUpdateTime.UpdateDefault.(func() time.Time)
	// guildconfigDescEnabled is the schema descriptor for enabled field.
	guildconfigDescEnabled := guildconfigFields[0].Descriptor()
	// guildconfig.DefaultEnabled holds the default value on creation for the enabled field.
	guildconfig.DefaultEnabled = guildconfigDescEnabled.Default.(bool)
	// guildconfigDescDefaultMaxClones is the schema descriptor for default_max_clones field.
	guildconfigDescDefaultMaxClones := guildconfigFields[1].Descriptor()
	// guildconfig.DefaultDefaultMaxClones holds the default value on creation for the default_max_clones field.
	guildconfig.DefaultDefaultMaxClones = guildconfigDescDefaultMaxClones.Default.(int)
	// guildconfig.DefaultMaxClonesValidator is a validator for the "default_max_clones" field. It is called by the builders before save.
	guildconfig.DefaultMaxClonesValidator = func() func(int) error {
		validators := guildconfigDescDefaultMaxClones.Validators
		fns := [...]func(int) error{
			validators[0].(func(int) error),
			validators[1].(func(int) error),
		}
		return func(default_max_clones int) error {
			for _, fn := range fns {
				if err := fn(default_max_clones); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// guildconfigDescRegexMatch is the schema descriptor for regex_match field.
	guildconfigDescRegexMatch := guildconfigFields[2].Descriptor()
	// guildconfig.DefaultRegexMatch holds the default value on creation for the regex_match field.
	guildconfig.DefaultRegexMatch = guildconfigDescRegexMatch.Default.(string)
	// guildconfig.RegexMatchValidator is a validator for the "regex_match" field. It is called by the builders before save.
	guildconfig.RegexMatchValidator = func() func(string) error {
		validators := guildconfigDescRegexMatch.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(regex_match string) error {
			for _, fn := range fns {
				if err := fn(regex_match); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// guildconfigDescContactEmail is the schema descriptor for contact_email field.
	guildconfigDescContactEmail := guildconfigFields[3].Descriptor()
	// guildconfig.DefaultContactEmail holds the default value on creation for the contact_email field.
	guildconfig.DefaultContactEmail = guildconfigDescContactEmail.Default.(string)
	// guildconfig.ContactEmailValidator is a validator for the "contact_email" field. It is called by the builders before save.
	guildconfig.ContactEmailValidator = func() func(string) error {
		validators := guildconfigDescContactEmail.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(contact_email string) error {
			for _, fn := range fns {
				if err := fn(contact_email); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	guildeventMixin := schema.GuildEvent{}.Mixin()
	guildevent.Policy = privacy.NewPolicies(schema.GuildEvent{})
	guildevent.Hooks[0] = func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if err := guildevent.Policy.EvalMutation(ctx, m); err != nil {
				return nil, err
			}
			return next.Mutate(ctx, m)
		})
	}
	guildeventMixinFields0 := guildeventMixin[0].Fields()
	_ = guildeventMixinFields0
	guildeventMixinFields1 := guildeventMixin[1].Fields()
	_ = guildeventMixinFields1
	guildeventFields := schema.GuildEvent{}.Fields()
	_ = guildeventFields
	// guildeventDescCreateTime is the schema descriptor for create_time field.
	guildeventDescCreateTime := guildeventMixinFields0[0].Descriptor()
	// guildevent.DefaultCreateTime holds the default value on creation for the create_time field.
	guildevent.DefaultCreateTime = guildeventDescCreateTime.Default.(func() time.Time)
	// guildeventDescUpdateTime is the schema descriptor for update_time field.
	guildeventDescUpdateTime := guildeventMixinFields1[0].Descriptor()
	// guildevent.DefaultUpdateTime holds the default value on creation for the update_time field.
	guildevent.DefaultUpdateTime = guildeventDescUpdateTime.Default.(func() time.Time)
	// guildevent.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	guildevent.UpdateDefaultUpdateTime = guildeventDescUpdateTime.UpdateDefault.(func() time.Time)
	userMixin := schema.User{}.Mixin()
	user.Policy = privacy.NewPolicies(schema.User{})
	user.Hooks[0] = func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if err := user.Policy.EvalMutation(ctx, m); err != nil {
				return nil, err
			}
			return next.Mutate(ctx, m)
		})
	}
	userMixinFields0 := userMixin[0].Fields()
	_ = userMixinFields0
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescCreateTime is the schema descriptor for create_time field.
	userDescCreateTime := userMixinFields0[0].Descriptor()
	// user.DefaultCreateTime holds the default value on creation for the create_time field.
	user.DefaultCreateTime = userDescCreateTime.Default.(func() time.Time)
	// userDescUpdateTime is the schema descriptor for update_time field.
	userDescUpdateTime := userMixinFields0[1].Descriptor()
	// user.DefaultUpdateTime holds the default value on creation for the update_time field.
	user.DefaultUpdateTime = userDescUpdateTime.Default.(func() time.Time)
	// user.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	user.UpdateDefaultUpdateTime = userDescUpdateTime.UpdateDefault.(func() time.Time)
	// userDescDiscriminator is the schema descriptor for discriminator field.
	userDescDiscriminator := userFields[5].Descriptor()
	// user.DiscriminatorValidator is a validator for the "discriminator" field. It is called by the builders before save.
	user.DiscriminatorValidator = userDescDiscriminator.Validators[0].(func(string) error)
	// userDescEmail is the schema descriptor for email field.
	userDescEmail := userFields[6].Descriptor()
	// user.EmailValidator is a validator for the "email" field. It is called by the builders before save.
	user.EmailValidator = userDescEmail.Validators[0].(func(string) error)
	// userDescAvatarHash is the schema descriptor for avatar_hash field.
	userDescAvatarHash := userFields[7].Descriptor()
	// user.AvatarHashValidator is a validator for the "avatar_hash" field. It is called by the builders before save.
	user.AvatarHashValidator = userDescAvatarHash.Validators[0].(func(string) error)
	// userDescAvatarURL is the schema descriptor for avatar_url field.
	userDescAvatarURL := userFields[8].Descriptor()
	// user.AvatarURLValidator is a validator for the "avatar_url" field. It is called by the builders before save.
	user.AvatarURLValidator = userDescAvatarURL.Validators[0].(func(string) error)
	// userDescLocale is the schema descriptor for locale field.
	userDescLocale := userFields[9].Descriptor()
	// user.LocaleValidator is a validator for the "locale" field. It is called by the builders before save.
	user.LocaleValidator = userDescLocale.Validators[0].(func(string) error)
}

const (
	Version = "v0.11.10"                                        // Version of ent codegen.
	Sum     = "h1:iqn32ybY5HRW3xSAyMNdNKpZhKgMf1Zunsej9yPKUI8=" // Sum of ent codegen.
)
