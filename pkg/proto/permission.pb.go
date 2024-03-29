// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.11.4
// source: pkg/proto/permission.proto

package proto

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// VerifyHttpRouteRequest 验证HTTP路由权限请求参数
type VerifyHttpRouteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// oauth2 client id
	Oauth2ClientId string `protobuf:"bytes,1,opt,name=oauth2_client_id,json=oauth2ClientId,proto3" json:"oauth2_client_id,omitempty"`
	// token
	Token string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	// 路由
	Path string `protobuf:"bytes,3,opt,name=path,proto3" json:"path,omitempty"`
	// 请求方法
	Method string `protobuf:"bytes,4,opt,name=method,proto3" json:"method,omitempty"`
	// 是否返回用户信息,token验证通过的情况下
	ReturnUserInfo bool `protobuf:"varint,5,opt,name=return_user_info,json=returnUserInfo,proto3" json:"return_user_info,omitempty"`
}

func (x *VerifyHttpRouteRequest) Reset() {
	*x = VerifyHttpRouteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_permission_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyHttpRouteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyHttpRouteRequest) ProtoMessage() {}

func (x *VerifyHttpRouteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_permission_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyHttpRouteRequest.ProtoReflect.Descriptor instead.
func (*VerifyHttpRouteRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_permission_proto_rawDescGZIP(), []int{0}
}

func (x *VerifyHttpRouteRequest) GetOauth2ClientId() string {
	if x != nil {
		return x.Oauth2ClientId
	}
	return ""
}

func (x *VerifyHttpRouteRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *VerifyHttpRouteRequest) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *VerifyHttpRouteRequest) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *VerifyHttpRouteRequest) GetReturnUserInfo() bool {
	if x != nil {
		return x.ReturnUserInfo
	}
	return false
}

// VerificationHttpRouterResponse 验证HTTP路由权限响应参数
type VerifyHttpRouteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 是否允许
	Allow bool `protobuf:"varint,1,opt,name=allow,proto3" json:"allow,omitempty"`
	// 用户信息
	UserInfo *VerifyHttpRouteResponse_UserInfo `protobuf:"bytes,2,opt,name=user_info,json=userInfo,proto3" json:"user_info,omitempty"`
}

func (x *VerifyHttpRouteResponse) Reset() {
	*x = VerifyHttpRouteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_permission_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyHttpRouteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyHttpRouteResponse) ProtoMessage() {}

func (x *VerifyHttpRouteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_permission_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyHttpRouteResponse.ProtoReflect.Descriptor instead.
func (*VerifyHttpRouteResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_permission_proto_rawDescGZIP(), []int{1}
}

func (x *VerifyHttpRouteResponse) GetAllow() bool {
	if x != nil {
		return x.Allow
	}
	return false
}

func (x *VerifyHttpRouteResponse) GetUserInfo() *VerifyHttpRouteResponse_UserInfo {
	if x != nil {
		return x.UserInfo
	}
	return nil
}

// VerificationTokenRequest 验证Token请求参数
type VerifyTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// oauth2 client id
	Oauth2ClientId string `protobuf:"bytes,1,opt,name=oauth2_client_id,json=oauth2ClientId,proto3" json:"oauth2_client_id,omitempty"`
	// token
	Token string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	// 是否返回用户信息,token验证通过的情况下
	ReturnUserInfo bool `protobuf:"varint,3,opt,name=return_user_info,json=returnUserInfo,proto3" json:"return_user_info,omitempty"`
}

func (x *VerifyTokenRequest) Reset() {
	*x = VerifyTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_permission_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyTokenRequest) ProtoMessage() {}

func (x *VerifyTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_permission_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyTokenRequest.ProtoReflect.Descriptor instead.
func (*VerifyTokenRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_permission_proto_rawDescGZIP(), []int{2}
}

func (x *VerifyTokenRequest) GetOauth2ClientId() string {
	if x != nil {
		return x.Oauth2ClientId
	}
	return ""
}

func (x *VerifyTokenRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *VerifyTokenRequest) GetReturnUserInfo() bool {
	if x != nil {
		return x.ReturnUserInfo
	}
	return false
}

// VerificationTokenResponse 验证Token响应参数
type VerifyTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 是否允许
	Allow bool `protobuf:"varint,1,opt,name=allow,proto3" json:"allow,omitempty"`
	// 用户信息
	UserInfo *VerifyTokenResponse_UserInfo `protobuf:"bytes,2,opt,name=user_info,json=userInfo,proto3" json:"user_info,omitempty"`
}

func (x *VerifyTokenResponse) Reset() {
	*x = VerifyTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_permission_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyTokenResponse) ProtoMessage() {}

func (x *VerifyTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_permission_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyTokenResponse.ProtoReflect.Descriptor instead.
func (*VerifyTokenResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_permission_proto_rawDescGZIP(), []int{3}
}

func (x *VerifyTokenResponse) GetAllow() bool {
	if x != nil {
		return x.Allow
	}
	return false
}

func (x *VerifyTokenResponse) GetUserInfo() *VerifyTokenResponse_UserInfo {
	if x != nil {
		return x.UserInfo
	}
	return nil
}

// 用户信息
type VerifyHttpRouteResponse_UserInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OpenId    string `protobuf:"bytes,1,opt,name=open_id,json=openId,proto3" json:"open_id,omitempty"`
	Username  string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	NickName  string `protobuf:"bytes,3,opt,name=nick_name,json=nickName,proto3" json:"nick_name,omitempty"`
	AvatarUrl string `protobuf:"bytes,4,opt,name=avatar_url,json=avatarUrl,proto3" json:"avatar_url,omitempty"`
	Gender    uint32 `protobuf:"varint,5,opt,name=gender,proto3" json:"gender,omitempty"`
}

func (x *VerifyHttpRouteResponse_UserInfo) Reset() {
	*x = VerifyHttpRouteResponse_UserInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_permission_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyHttpRouteResponse_UserInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyHttpRouteResponse_UserInfo) ProtoMessage() {}

func (x *VerifyHttpRouteResponse_UserInfo) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_permission_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyHttpRouteResponse_UserInfo.ProtoReflect.Descriptor instead.
func (*VerifyHttpRouteResponse_UserInfo) Descriptor() ([]byte, []int) {
	return file_pkg_proto_permission_proto_rawDescGZIP(), []int{1, 0}
}

func (x *VerifyHttpRouteResponse_UserInfo) GetOpenId() string {
	if x != nil {
		return x.OpenId
	}
	return ""
}

func (x *VerifyHttpRouteResponse_UserInfo) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *VerifyHttpRouteResponse_UserInfo) GetNickName() string {
	if x != nil {
		return x.NickName
	}
	return ""
}

func (x *VerifyHttpRouteResponse_UserInfo) GetAvatarUrl() string {
	if x != nil {
		return x.AvatarUrl
	}
	return ""
}

func (x *VerifyHttpRouteResponse_UserInfo) GetGender() uint32 {
	if x != nil {
		return x.Gender
	}
	return 0
}

// 用户信息
type VerifyTokenResponse_UserInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OpenId    string `protobuf:"bytes,1,opt,name=open_id,json=openId,proto3" json:"open_id,omitempty"`
	Username  string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	NickName  string `protobuf:"bytes,3,opt,name=nick_name,json=nickName,proto3" json:"nick_name,omitempty"`
	AvatarUrl string `protobuf:"bytes,4,opt,name=avatar_url,json=avatarUrl,proto3" json:"avatar_url,omitempty"`
	Gender    uint32 `protobuf:"varint,5,opt,name=gender,proto3" json:"gender,omitempty"`
}

func (x *VerifyTokenResponse_UserInfo) Reset() {
	*x = VerifyTokenResponse_UserInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_permission_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyTokenResponse_UserInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyTokenResponse_UserInfo) ProtoMessage() {}

func (x *VerifyTokenResponse_UserInfo) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_permission_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyTokenResponse_UserInfo.ProtoReflect.Descriptor instead.
func (*VerifyTokenResponse_UserInfo) Descriptor() ([]byte, []int) {
	return file_pkg_proto_permission_proto_rawDescGZIP(), []int{3, 0}
}

func (x *VerifyTokenResponse_UserInfo) GetOpenId() string {
	if x != nil {
		return x.OpenId
	}
	return ""
}

func (x *VerifyTokenResponse_UserInfo) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *VerifyTokenResponse_UserInfo) GetNickName() string {
	if x != nil {
		return x.NickName
	}
	return ""
}

func (x *VerifyTokenResponse_UserInfo) GetAvatarUrl() string {
	if x != nil {
		return x.AvatarUrl
	}
	return ""
}

func (x *VerifyTokenResponse_UserInfo) GetGender() uint32 {
	if x != nil {
		return x.Gender
	}
	return 0
}

var File_pkg_proto_permission_proto protoreflect.FileDescriptor

var file_pkg_proto_permission_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x6e, 0x69,
	0x6c, 0x6f, 0x72, 0x67, 0x2e, 0x6e, 0x61, 0x61, 0x73, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xae, 0x01, 0x0a, 0x16, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x48, 0x74, 0x74, 0x70,
	0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x10,
	0x6f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x43, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x12, 0x0a, 0x04,
	0x70, 0x61, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68,
	0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x28, 0x0a, 0x10, 0x72, 0x65, 0x74, 0x75,
	0x72, 0x6e, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x0e, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e,
	0x66, 0x6f, 0x22, 0x9b, 0x02, 0x0a, 0x17, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x48, 0x74, 0x74,
	0x70, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x61,
	0x6c, 0x6c, 0x6f, 0x77, 0x12, 0x54, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x66,
	0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x37, 0x2e, 0x6e, 0x69, 0x6c, 0x6f, 0x72, 0x67,
	0x2e, 0x6e, 0x61, 0x61, 0x73, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x48, 0x74, 0x74, 0x70, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x93, 0x01, 0x0a, 0x08, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x17, 0x0a, 0x07, 0x6f, 0x70, 0x65, 0x6e, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x70, 0x65, 0x6e, 0x49, 0x64,
	0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x6e, 0x69, 0x63, 0x6b, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x6e, 0x69, 0x63, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x76, 0x61,
	0x74, 0x61, 0x72, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61,
	0x76, 0x61, 0x74, 0x61, 0x72, 0x55, 0x72, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x22, 0x7e, 0x0a, 0x12, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x10, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x32,
	0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0e, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x28, 0x0a, 0x10, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e,
	0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0e, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f,
	0x22, 0x93, 0x02, 0x0a, 0x13, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x6c, 0x6f,
	0x77, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x12, 0x50,
	0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x33, 0x2e, 0x6e, 0x69, 0x6c, 0x6f, 0x72, 0x67, 0x2e, 0x6e, 0x61, 0x61, 0x73, 0x2e,
	0x70, 0x6b, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f,
	0x1a, 0x93, 0x01, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x17, 0x0a,
	0x07, 0x6f, 0x70, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x6f, 0x70, 0x65, 0x6e, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x69, 0x63, 0x6b, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x1d, 0x0a, 0x0a, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x55, 0x72, 0x6c, 0x12, 0x16,
	0x0a, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06,
	0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x32, 0xc2, 0x02, 0x0a, 0x0a, 0x50, 0x65, 0x72, 0x6d, 0x69,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x8f, 0x01, 0x0a, 0x0b, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x29, 0x2e, 0x6e, 0x69, 0x6c, 0x6f, 0x72, 0x67, 0x2e, 0x6e,
	0x61, 0x61, 0x73, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x56, 0x65,
	0x72, 0x69, 0x66, 0x79, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x2a, 0x2e, 0x6e, 0x69, 0x6c, 0x6f, 0x72, 0x67, 0x2e, 0x6e, 0x61, 0x61, 0x73, 0x2e, 0x70,
	0x6b, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x29, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x23, 0x22, 0x1e, 0x2f, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x2f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x3a, 0x01, 0x2a, 0x12, 0xa1, 0x01, 0x0a, 0x0f, 0x56, 0x65, 0x72, 0x69,
	0x66, 0x79, 0x48, 0x74, 0x74, 0x70, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x12, 0x2d, 0x2e, 0x6e, 0x69,
	0x6c, 0x6f, 0x72, 0x67, 0x2e, 0x6e, 0x61, 0x61, 0x73, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x48, 0x74, 0x74, 0x70, 0x52, 0x6f,
	0x75, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2e, 0x2e, 0x6e, 0x69, 0x6c,
	0x6f, 0x72, 0x67, 0x2e, 0x6e, 0x61, 0x61, 0x73, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x48, 0x74, 0x74, 0x70, 0x52, 0x6f, 0x75,
	0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2f, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x29, 0x22, 0x24, 0x2f, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2f,
	0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x68, 0x74, 0x74,
	0x70, 0x5f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x3a, 0x01, 0x2a, 0x42, 0x22, 0x5a, 0x20, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x69, 0x6c, 0x6f, 0x72, 0x67,
	0x2f, 0x6e, 0x61, 0x61, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_proto_permission_proto_rawDescOnce sync.Once
	file_pkg_proto_permission_proto_rawDescData = file_pkg_proto_permission_proto_rawDesc
)

func file_pkg_proto_permission_proto_rawDescGZIP() []byte {
	file_pkg_proto_permission_proto_rawDescOnce.Do(func() {
		file_pkg_proto_permission_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_proto_permission_proto_rawDescData)
	})
	return file_pkg_proto_permission_proto_rawDescData
}

var file_pkg_proto_permission_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_pkg_proto_permission_proto_goTypes = []interface{}{
	(*VerifyHttpRouteRequest)(nil),           // 0: nilorg.naas.pkg.proto.VerifyHttpRouteRequest
	(*VerifyHttpRouteResponse)(nil),          // 1: nilorg.naas.pkg.proto.VerifyHttpRouteResponse
	(*VerifyTokenRequest)(nil),               // 2: nilorg.naas.pkg.proto.VerifyTokenRequest
	(*VerifyTokenResponse)(nil),              // 3: nilorg.naas.pkg.proto.VerifyTokenResponse
	(*VerifyHttpRouteResponse_UserInfo)(nil), // 4: nilorg.naas.pkg.proto.VerifyHttpRouteResponse.UserInfo
	(*VerifyTokenResponse_UserInfo)(nil),     // 5: nilorg.naas.pkg.proto.VerifyTokenResponse.UserInfo
}
var file_pkg_proto_permission_proto_depIdxs = []int32{
	4, // 0: nilorg.naas.pkg.proto.VerifyHttpRouteResponse.user_info:type_name -> nilorg.naas.pkg.proto.VerifyHttpRouteResponse.UserInfo
	5, // 1: nilorg.naas.pkg.proto.VerifyTokenResponse.user_info:type_name -> nilorg.naas.pkg.proto.VerifyTokenResponse.UserInfo
	2, // 2: nilorg.naas.pkg.proto.Permission.VerifyToken:input_type -> nilorg.naas.pkg.proto.VerifyTokenRequest
	0, // 3: nilorg.naas.pkg.proto.Permission.VerifyHttpRoute:input_type -> nilorg.naas.pkg.proto.VerifyHttpRouteRequest
	3, // 4: nilorg.naas.pkg.proto.Permission.VerifyToken:output_type -> nilorg.naas.pkg.proto.VerifyTokenResponse
	1, // 5: nilorg.naas.pkg.proto.Permission.VerifyHttpRoute:output_type -> nilorg.naas.pkg.proto.VerifyHttpRouteResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_pkg_proto_permission_proto_init() }
func file_pkg_proto_permission_proto_init() {
	if File_pkg_proto_permission_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_proto_permission_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyHttpRouteRequest); i {
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
		file_pkg_proto_permission_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyHttpRouteResponse); i {
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
		file_pkg_proto_permission_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyTokenRequest); i {
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
		file_pkg_proto_permission_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyTokenResponse); i {
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
		file_pkg_proto_permission_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyHttpRouteResponse_UserInfo); i {
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
		file_pkg_proto_permission_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyTokenResponse_UserInfo); i {
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
			RawDescriptor: file_pkg_proto_permission_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_proto_permission_proto_goTypes,
		DependencyIndexes: file_pkg_proto_permission_proto_depIdxs,
		MessageInfos:      file_pkg_proto_permission_proto_msgTypes,
	}.Build()
	File_pkg_proto_permission_proto = out.File
	file_pkg_proto_permission_proto_rawDesc = nil
	file_pkg_proto_permission_proto_goTypes = nil
	file_pkg_proto_permission_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PermissionClient is the client API for Permission service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PermissionClient interface {
	// VerifyToken 验证Token
	VerifyToken(ctx context.Context, in *VerifyTokenRequest, opts ...grpc.CallOption) (*VerifyTokenResponse, error)
	// VerifyHttpRoute 验证Http路由权限
	VerifyHttpRoute(ctx context.Context, in *VerifyHttpRouteRequest, opts ...grpc.CallOption) (*VerifyHttpRouteResponse, error)
}

type permissionClient struct {
	cc grpc.ClientConnInterface
}

func NewPermissionClient(cc grpc.ClientConnInterface) PermissionClient {
	return &permissionClient{cc}
}

func (c *permissionClient) VerifyToken(ctx context.Context, in *VerifyTokenRequest, opts ...grpc.CallOption) (*VerifyTokenResponse, error) {
	out := new(VerifyTokenResponse)
	err := c.cc.Invoke(ctx, "/nilorg.naas.pkg.proto.Permission/VerifyToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionClient) VerifyHttpRoute(ctx context.Context, in *VerifyHttpRouteRequest, opts ...grpc.CallOption) (*VerifyHttpRouteResponse, error) {
	out := new(VerifyHttpRouteResponse)
	err := c.cc.Invoke(ctx, "/nilorg.naas.pkg.proto.Permission/VerifyHttpRoute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PermissionServer is the server API for Permission service.
type PermissionServer interface {
	// VerifyToken 验证Token
	VerifyToken(context.Context, *VerifyTokenRequest) (*VerifyTokenResponse, error)
	// VerifyHttpRoute 验证Http路由权限
	VerifyHttpRoute(context.Context, *VerifyHttpRouteRequest) (*VerifyHttpRouteResponse, error)
}

// UnimplementedPermissionServer can be embedded to have forward compatible implementations.
type UnimplementedPermissionServer struct {
}

func (*UnimplementedPermissionServer) VerifyToken(context.Context, *VerifyTokenRequest) (*VerifyTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyToken not implemented")
}
func (*UnimplementedPermissionServer) VerifyHttpRoute(context.Context, *VerifyHttpRouteRequest) (*VerifyHttpRouteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyHttpRoute not implemented")
}

func RegisterPermissionServer(s *grpc.Server, srv PermissionServer) {
	s.RegisterService(&_Permission_serviceDesc, srv)
}

func _Permission_VerifyToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServer).VerifyToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nilorg.naas.pkg.proto.Permission/VerifyToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServer).VerifyToken(ctx, req.(*VerifyTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Permission_VerifyHttpRoute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyHttpRouteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServer).VerifyHttpRoute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nilorg.naas.pkg.proto.Permission/VerifyHttpRoute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServer).VerifyHttpRoute(ctx, req.(*VerifyHttpRouteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Permission_serviceDesc = grpc.ServiceDesc{
	ServiceName: "nilorg.naas.pkg.proto.Permission",
	HandlerType: (*PermissionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "VerifyToken",
			Handler:    _Permission_VerifyToken_Handler,
		},
		{
			MethodName: "VerifyHttpRoute",
			Handler:    _Permission_VerifyHttpRoute_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/proto/permission.proto",
}
