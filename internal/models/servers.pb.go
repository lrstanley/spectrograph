// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.15.7
// source: servers.proto

package models

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Server struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id" bson:"_id"`           // @inject_tag: bson:"_id" json:"id"
	Created *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=created,proto3" json:"created" bson:"created"` // @inject_tag: bson:"created" json:"created"
	Updated *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=updated,proto3" json:"updated" bson:"updated"` // @inject_tag: bson:"updated" json:"updated"
	Discord *ServerDiscordData     `protobuf:"bytes,4,opt,name=discord,proto3" json:"discord" bson:"discord"` // @inject_tag: bson:"discord" json:"discord"`
}

func (x *Server) Reset() {
	*x = Server{}
	if protoimpl.UnsafeEnabled {
		mi := &file_servers_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Server) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Server) ProtoMessage() {}

func (x *Server) ProtoReflect() protoreflect.Message {
	mi := &file_servers_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Server.ProtoReflect.Descriptor instead.
func (*Server) Descriptor() ([]byte, []int) {
	return file_servers_proto_rawDescGZIP(), []int{0}
}

func (x *Server) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Server) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

func (x *Server) GetUpdated() *timestamppb.Timestamp {
	if x != nil {
		return x.Updated
	}
	return nil
}

func (x *Server) GetDiscord() *ServerDiscordData {
	if x != nil {
		return x.Discord
	}
	return nil
}

type ServerDiscordData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Guild ID.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // @inject_tag: bson:"id" json:"id"
	// Guild name (2-100 chars, excl. trailing/leading spaces).
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"` // @inject_tag: bson:"name" json:"name"
	// Guild description.
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"` // @inject_tag: bson:"description" json:"description"
	// Enabled guild features.
	Features []string `protobuf:"bytes,4,rep,name=features,proto3" json:"features,omitempty"` // @inject_tag: bson:"features" json:"features"
	// Icon hash.
	Icon string `protobuf:"bytes,5,opt,name=icon,proto3" json:"icon,omitempty"` // @inject_tag: bson:"icon" json:"icon"
	// This is something we generate.
	IconUrl string `protobuf:"bytes,6,opt,name=icon_url,json=iconUrl,proto3" json:"icon_url,omitempty"` // @inject_tag: bson:"icon_url" json:"icon_url"
	// When the bot joined the guild.
	JoinedAt *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=joined_at,json=joinedAt,proto3" json:"joined_at,omitempty"` // @inject_tag: bson:"joined_at" json:"joined_at"
	// If the guild is considered large (to Discord standards).
	Large bool `protobuf:"varint,8,opt,name=large,proto3" json:"large,omitempty"` // @inject_tag: bson:"large" json:"large"
	// Total members in this guild.
	MemberCount int64 `protobuf:"varint,9,opt,name=member_count,json=memberCount,proto3" json:"member_count,omitempty"` // @inject_tag: bson:"member_count" json:"member_count"
	// User ID of the owner.
	OwnerId string `protobuf:"bytes,10,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"` // @inject_tag: bson:"owner_id" json:"owner_id"
	// Permissions of the bot on the server.
	Permissions uint64 `protobuf:"varint,11,opt,name=permissions,proto3" json:"permissions,omitempty"` // @inject_tag: bson:"permissions" json:"permissions"
	// Preferred locale.
	PreferredLocale string `protobuf:"bytes,12,opt,name=preferred_locale,json=preferredLocale,proto3" json:"preferred_locale,omitempty"` // @inject_tag: bson:"preferred_locale" json:"preferred_locale"
	// Voice region (deprecated?).
	Region                 string `protobuf:"bytes,13,opt,name=Region,proto3" json:"Region,omitempty"`                                                                   // @inject_tag: bson:"region" json:"region"
	PublicUpdatesChannelId string `protobuf:"bytes,14,opt,name=public_updates_channel_id,json=publicUpdatesChannelId,proto3" json:"public_updates_channel_id" bson:"public_updates_channel_id"` // @inject_tag: bson:"public_updates_channel_id" json:"public_updates_channel_id"
	SystemChannelFlags     string `protobuf:"bytes,15,opt,name=system_channel_flags,json=systemChannelFlags,proto3" json:"system_channel_flags" bson:"system_channel_flags"`               // @inject_tag: bson:"system_channel_flags" json:"system_channel_flags"
}

func (x *ServerDiscordData) Reset() {
	*x = ServerDiscordData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_servers_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerDiscordData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerDiscordData) ProtoMessage() {}

func (x *ServerDiscordData) ProtoReflect() protoreflect.Message {
	mi := &file_servers_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerDiscordData.ProtoReflect.Descriptor instead.
func (*ServerDiscordData) Descriptor() ([]byte, []int) {
	return file_servers_proto_rawDescGZIP(), []int{1}
}

func (x *ServerDiscordData) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ServerDiscordData) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ServerDiscordData) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ServerDiscordData) GetFeatures() []string {
	if x != nil {
		return x.Features
	}
	return nil
}

func (x *ServerDiscordData) GetIcon() string {
	if x != nil {
		return x.Icon
	}
	return ""
}

func (x *ServerDiscordData) GetIconUrl() string {
	if x != nil {
		return x.IconUrl
	}
	return ""
}

func (x *ServerDiscordData) GetJoinedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.JoinedAt
	}
	return nil
}

func (x *ServerDiscordData) GetLarge() bool {
	if x != nil {
		return x.Large
	}
	return false
}

func (x *ServerDiscordData) GetMemberCount() int64 {
	if x != nil {
		return x.MemberCount
	}
	return 0
}

func (x *ServerDiscordData) GetOwnerId() string {
	if x != nil {
		return x.OwnerId
	}
	return ""
}

func (x *ServerDiscordData) GetPermissions() uint64 {
	if x != nil {
		return x.Permissions
	}
	return 0
}

func (x *ServerDiscordData) GetPreferredLocale() string {
	if x != nil {
		return x.PreferredLocale
	}
	return ""
}

func (x *ServerDiscordData) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *ServerDiscordData) GetPublicUpdatesChannelId() string {
	if x != nil {
		return x.PublicUpdatesChannelId
	}
	return ""
}

func (x *ServerDiscordData) GetSystemChannelFlags() string {
	if x != nil {
		return x.SystemChannelFlags
	}
	return ""
}

// TODO: auto-generate status if last status message is greater than X period
// of time?
type ServerStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// string id = 1; // @inject_tag: bson:"_id" json:"id"
	Id        string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id" bson:"_id"`                // @gotags:[bson:"_id" json:"id"]
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp" bson:"timestamp"`  // @inject_tag: bson:"timestamp" json:"timestamp"
	Type      string                 `protobuf:"bytes,3,opt,name=type,proto3" json:"type" bson:"type"`            // @inject_tag: bson:"type" json:"type"
	Message   string                 `protobuf:"bytes,4,opt,name=message,proto3" json:"message" bson:"message"`      // @inject_tag: bson:"message" json:"message"
	Healthy   bool                   `protobuf:"varint,5,opt,name=healthy,proto3" json:"healthy"`     // @inject_tag: bson:'healthy" json:"healthy"
	Available bool                   `protobuf:"varint,6,opt,name=available,proto3" json:"available" bson:"available"` // @inject_tag: bson:"available" json:"available"
}

func (x *ServerStatus) Reset() {
	*x = ServerStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_servers_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerStatus) ProtoMessage() {}

func (x *ServerStatus) ProtoReflect() protoreflect.Message {
	mi := &file_servers_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerStatus.ProtoReflect.Descriptor instead.
func (*ServerStatus) Descriptor() ([]byte, []int) {
	return file_servers_proto_rawDescGZIP(), []int{2}
}

func (x *ServerStatus) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ServerStatus) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *ServerStatus) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *ServerStatus) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *ServerStatus) GetHealthy() bool {
	if x != nil {
		return x.Healthy
	}
	return false
}

func (x *ServerStatus) GetAvailable() bool {
	if x != nil {
		return x.Available
	}
	return false
}

var File_servers_proto protoreflect.FileDescriptor

var file_servers_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x06, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb9, 0x01, 0x0a, 0x06, 0x53, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x34, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x34, 0x0a, 0x07, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x12,
	0x33, 0x0a, 0x07, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x44, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x44, 0x61, 0x74, 0x61, 0x52, 0x07, 0x64, 0x69, 0x73,
	0x63, 0x6f, 0x72, 0x64, 0x22, 0x83, 0x04, 0x0a, 0x11, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x44,
	0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x44, 0x61, 0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x1a, 0x0a, 0x08, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x08, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x69, 0x63, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x63, 0x6f, 0x6e,
	0x12, 0x19, 0x0a, 0x08, 0x69, 0x63, 0x6f, 0x6e, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x69, 0x63, 0x6f, 0x6e, 0x55, 0x72, 0x6c, 0x12, 0x37, 0x0a, 0x09, 0x6a,
	0x6f, 0x69, 0x6e, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x6a, 0x6f, 0x69, 0x6e,
	0x65, 0x64, 0x41, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x61, 0x72, 0x67, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x05, 0x6c, 0x61, 0x72, 0x67, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x65,
	0x6d, 0x62, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0b, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x19, 0x0a,
	0x08, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x70,
	0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x29, 0x0a, 0x10, 0x70, 0x72,
	0x65, 0x66, 0x65, 0x72, 0x72, 0x65, 0x64, 0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x65, 0x18, 0x0c,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x70, 0x72, 0x65, 0x66, 0x65, 0x72, 0x72, 0x65, 0x64, 0x4c,
	0x6f, 0x63, 0x61, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18,
	0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x39, 0x0a,
	0x19, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x5f,
	0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x16, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x43,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x64, 0x12, 0x30, 0x0a, 0x14, 0x73, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x66, 0x6c, 0x61, 0x67, 0x73,
	0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x43, 0x68,
	0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x46, 0x6c, 0x61, 0x67, 0x73, 0x22, 0xbe, 0x01, 0x0a, 0x0c, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x38, 0x0a, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x79, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x79, 0x12, 0x1c, 0x0a,
	0x09, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x09, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x42, 0x33, 0x5a, 0x31, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x72, 0x73, 0x74, 0x61, 0x6e,
	0x6c, 0x65, 0x79, 0x2f, 0x73, 0x70, 0x65, 0x63, 0x74, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x70, 0x68,
	0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_servers_proto_rawDescOnce sync.Once
	file_servers_proto_rawDescData = file_servers_proto_rawDesc
)

func file_servers_proto_rawDescGZIP() []byte {
	file_servers_proto_rawDescOnce.Do(func() {
		file_servers_proto_rawDescData = protoimpl.X.CompressGZIP(file_servers_proto_rawDescData)
	})
	return file_servers_proto_rawDescData
}

var file_servers_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_servers_proto_goTypes = []interface{}{
	(*Server)(nil),                // 0: models.Server
	(*ServerDiscordData)(nil),     // 1: models.ServerDiscordData
	(*ServerStatus)(nil),          // 2: models.ServerStatus
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_servers_proto_depIdxs = []int32{
	3, // 0: models.Server.created:type_name -> google.protobuf.Timestamp
	3, // 1: models.Server.updated:type_name -> google.protobuf.Timestamp
	1, // 2: models.Server.discord:type_name -> models.ServerDiscordData
	3, // 3: models.ServerDiscordData.joined_at:type_name -> google.protobuf.Timestamp
	3, // 4: models.ServerStatus.timestamp:type_name -> google.protobuf.Timestamp
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_servers_proto_init() }
func file_servers_proto_init() {
	if File_servers_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_servers_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Server); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_servers_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerDiscordData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_servers_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerStatus); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_servers_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_servers_proto_goTypes,
		DependencyIndexes: file_servers_proto_depIdxs,
		MessageInfos:      file_servers_proto_msgTypes,
	}.Build()
	File_servers_proto = out.File
	file_servers_proto_rawDesc = nil
	file_servers_proto_goTypes = nil
	file_servers_proto_depIdxs = nil
}
