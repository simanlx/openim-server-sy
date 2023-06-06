// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: user/user.proto

package user

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	User_GetUserInfo_FullMethodName             = "/user.user/GetUserInfo"
	User_UpdateUserInfo_FullMethodName          = "/user.user/UpdateUserInfo"
	User_SetGlobalRecvMessageOpt_FullMethodName = "/user.user/SetGlobalRecvMessageOpt"
	User_GetAllUserID_FullMethodName            = "/user.user/GetAllUserID"
	User_AccountCheck_FullMethodName            = "/user.user/AccountCheck"
	User_GetConversation_FullMethodName         = "/user.user/GetConversation"
	User_GetAllConversations_FullMethodName     = "/user.user/GetAllConversations"
	User_GetConversations_FullMethodName        = "/user.user/GetConversations"
	User_BatchSetConversations_FullMethodName   = "/user.user/BatchSetConversations"
	User_SetConversation_FullMethodName         = "/user.user/SetConversation"
	User_SetRecvMsgOpt_FullMethodName           = "/user.user/SetRecvMsgOpt"
	User_GetUsers_FullMethodName                = "/user.user/GetUsers"
	User_AddUser_FullMethodName                 = "/user.user/AddUser"
	User_BlockUser_FullMethodName               = "/user.user/BlockUser"
	User_UnBlockUser_FullMethodName             = "/user.user/UnBlockUser"
	User_GetBlockUsers_FullMethodName           = "/user.user/GetBlockUsers"
	User_AttributeSwitch_FullMethodName         = "/user.user/AttributeSwitch"
	User_AttributeSwitchSet_FullMethodName      = "/user.user/AttributeSwitchSet"
	User_AttributeMenu_FullMethodName           = "/user.user/AttributeMenu"
	User_Feedback_FullMethodName                = "/user.user/Feedback"
	User_GetCommonProblem_FullMethodName        = "/user.user/GetCommonProblem"
	User_FeedbackCommonProblem_FullMethodName   = "/user.user/FeedbackCommonProblem"
)

// UserClient is the client API for User service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserClient interface {
	GetUserInfo(ctx context.Context, in *GetUserInfoReq, opts ...grpc.CallOption) (*GetUserInfoResp, error)
	UpdateUserInfo(ctx context.Context, in *UpdateUserInfoReq, opts ...grpc.CallOption) (*UpdateUserInfoResp, error)
	SetGlobalRecvMessageOpt(ctx context.Context, in *SetGlobalRecvMessageOptReq, opts ...grpc.CallOption) (*SetGlobalRecvMessageOptResp, error)
	GetAllUserID(ctx context.Context, in *GetAllUserIDReq, opts ...grpc.CallOption) (*GetAllUserIDResp, error)
	AccountCheck(ctx context.Context, in *AccountCheckReq, opts ...grpc.CallOption) (*AccountCheckResp, error)
	GetConversation(ctx context.Context, in *GetConversationReq, opts ...grpc.CallOption) (*GetConversationResp, error)
	GetAllConversations(ctx context.Context, in *GetAllConversationsReq, opts ...grpc.CallOption) (*GetAllConversationsResp, error)
	GetConversations(ctx context.Context, in *GetConversationsReq, opts ...grpc.CallOption) (*GetConversationsResp, error)
	BatchSetConversations(ctx context.Context, in *BatchSetConversationsReq, opts ...grpc.CallOption) (*BatchSetConversationsResp, error)
	SetConversation(ctx context.Context, in *SetConversationReq, opts ...grpc.CallOption) (*SetConversationResp, error)
	SetRecvMsgOpt(ctx context.Context, in *SetRecvMsgOptReq, opts ...grpc.CallOption) (*SetRecvMsgOptResp, error)
	GetUsers(ctx context.Context, in *GetUsersReq, opts ...grpc.CallOption) (*GetUsersResp, error)
	AddUser(ctx context.Context, in *AddUserReq, opts ...grpc.CallOption) (*AddUserResp, error)
	BlockUser(ctx context.Context, in *BlockUserReq, opts ...grpc.CallOption) (*BlockUserResp, error)
	UnBlockUser(ctx context.Context, in *UnBlockUserReq, opts ...grpc.CallOption) (*UnBlockUserResp, error)
	GetBlockUsers(ctx context.Context, in *GetBlockUsersReq, opts ...grpc.CallOption) (*GetBlockUsersResp, error)
	AttributeSwitch(ctx context.Context, in *AttributeSwitchReq, opts ...grpc.CallOption) (*AttributeSwitchResp, error)
	AttributeSwitchSet(ctx context.Context, in *AttributeSwitchSetReq, opts ...grpc.CallOption) (*AttributeSwitchSetResp, error)
	AttributeMenu(ctx context.Context, in *AttributeMenuReq, opts ...grpc.CallOption) (*AttributeMenuResp, error)
	// 用户反馈
	Feedback(ctx context.Context, in *FeedbackReq, opts ...grpc.CallOption) (*FeedbackResp, error)
	// 常见问题
	GetCommonProblem(ctx context.Context, in *GetCommonProblemReq, opts ...grpc.CallOption) (*GetCommonProblemResp, error)
	// 常见问题-反馈
	FeedbackCommonProblem(ctx context.Context, in *FeedbackCommonProblemReq, opts ...grpc.CallOption) (*FeedbackCommonProblemResp, error)
}

type userClient struct {
	cc grpc.ClientConnInterface
}

func NewUserClient(cc grpc.ClientConnInterface) UserClient {
	return &userClient{cc}
}

func (c *userClient) GetUserInfo(ctx context.Context, in *GetUserInfoReq, opts ...grpc.CallOption) (*GetUserInfoResp, error) {
	out := new(GetUserInfoResp)
	err := c.cc.Invoke(ctx, User_GetUserInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) UpdateUserInfo(ctx context.Context, in *UpdateUserInfoReq, opts ...grpc.CallOption) (*UpdateUserInfoResp, error) {
	out := new(UpdateUserInfoResp)
	err := c.cc.Invoke(ctx, User_UpdateUserInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) SetGlobalRecvMessageOpt(ctx context.Context, in *SetGlobalRecvMessageOptReq, opts ...grpc.CallOption) (*SetGlobalRecvMessageOptResp, error) {
	out := new(SetGlobalRecvMessageOptResp)
	err := c.cc.Invoke(ctx, User_SetGlobalRecvMessageOpt_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetAllUserID(ctx context.Context, in *GetAllUserIDReq, opts ...grpc.CallOption) (*GetAllUserIDResp, error) {
	out := new(GetAllUserIDResp)
	err := c.cc.Invoke(ctx, User_GetAllUserID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) AccountCheck(ctx context.Context, in *AccountCheckReq, opts ...grpc.CallOption) (*AccountCheckResp, error) {
	out := new(AccountCheckResp)
	err := c.cc.Invoke(ctx, User_AccountCheck_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetConversation(ctx context.Context, in *GetConversationReq, opts ...grpc.CallOption) (*GetConversationResp, error) {
	out := new(GetConversationResp)
	err := c.cc.Invoke(ctx, User_GetConversation_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetAllConversations(ctx context.Context, in *GetAllConversationsReq, opts ...grpc.CallOption) (*GetAllConversationsResp, error) {
	out := new(GetAllConversationsResp)
	err := c.cc.Invoke(ctx, User_GetAllConversations_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetConversations(ctx context.Context, in *GetConversationsReq, opts ...grpc.CallOption) (*GetConversationsResp, error) {
	out := new(GetConversationsResp)
	err := c.cc.Invoke(ctx, User_GetConversations_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) BatchSetConversations(ctx context.Context, in *BatchSetConversationsReq, opts ...grpc.CallOption) (*BatchSetConversationsResp, error) {
	out := new(BatchSetConversationsResp)
	err := c.cc.Invoke(ctx, User_BatchSetConversations_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) SetConversation(ctx context.Context, in *SetConversationReq, opts ...grpc.CallOption) (*SetConversationResp, error) {
	out := new(SetConversationResp)
	err := c.cc.Invoke(ctx, User_SetConversation_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) SetRecvMsgOpt(ctx context.Context, in *SetRecvMsgOptReq, opts ...grpc.CallOption) (*SetRecvMsgOptResp, error) {
	out := new(SetRecvMsgOptResp)
	err := c.cc.Invoke(ctx, User_SetRecvMsgOpt_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetUsers(ctx context.Context, in *GetUsersReq, opts ...grpc.CallOption) (*GetUsersResp, error) {
	out := new(GetUsersResp)
	err := c.cc.Invoke(ctx, User_GetUsers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) AddUser(ctx context.Context, in *AddUserReq, opts ...grpc.CallOption) (*AddUserResp, error) {
	out := new(AddUserResp)
	err := c.cc.Invoke(ctx, User_AddUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) BlockUser(ctx context.Context, in *BlockUserReq, opts ...grpc.CallOption) (*BlockUserResp, error) {
	out := new(BlockUserResp)
	err := c.cc.Invoke(ctx, User_BlockUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) UnBlockUser(ctx context.Context, in *UnBlockUserReq, opts ...grpc.CallOption) (*UnBlockUserResp, error) {
	out := new(UnBlockUserResp)
	err := c.cc.Invoke(ctx, User_UnBlockUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetBlockUsers(ctx context.Context, in *GetBlockUsersReq, opts ...grpc.CallOption) (*GetBlockUsersResp, error) {
	out := new(GetBlockUsersResp)
	err := c.cc.Invoke(ctx, User_GetBlockUsers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) AttributeSwitch(ctx context.Context, in *AttributeSwitchReq, opts ...grpc.CallOption) (*AttributeSwitchResp, error) {
	out := new(AttributeSwitchResp)
	err := c.cc.Invoke(ctx, User_AttributeSwitch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) AttributeSwitchSet(ctx context.Context, in *AttributeSwitchSetReq, opts ...grpc.CallOption) (*AttributeSwitchSetResp, error) {
	out := new(AttributeSwitchSetResp)
	err := c.cc.Invoke(ctx, User_AttributeSwitchSet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) AttributeMenu(ctx context.Context, in *AttributeMenuReq, opts ...grpc.CallOption) (*AttributeMenuResp, error) {
	out := new(AttributeMenuResp)
	err := c.cc.Invoke(ctx, User_AttributeMenu_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) Feedback(ctx context.Context, in *FeedbackReq, opts ...grpc.CallOption) (*FeedbackResp, error) {
	out := new(FeedbackResp)
	err := c.cc.Invoke(ctx, User_Feedback_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetCommonProblem(ctx context.Context, in *GetCommonProblemReq, opts ...grpc.CallOption) (*GetCommonProblemResp, error) {
	out := new(GetCommonProblemResp)
	err := c.cc.Invoke(ctx, User_GetCommonProblem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) FeedbackCommonProblem(ctx context.Context, in *FeedbackCommonProblemReq, opts ...grpc.CallOption) (*FeedbackCommonProblemResp, error) {
	out := new(FeedbackCommonProblemResp)
	err := c.cc.Invoke(ctx, User_FeedbackCommonProblem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServer is the server API for User service.
// All implementations must embed UnimplementedUserServer
// for forward compatibility
type UserServer interface {
	GetUserInfo(context.Context, *GetUserInfoReq) (*GetUserInfoResp, error)
	UpdateUserInfo(context.Context, *UpdateUserInfoReq) (*UpdateUserInfoResp, error)
	SetGlobalRecvMessageOpt(context.Context, *SetGlobalRecvMessageOptReq) (*SetGlobalRecvMessageOptResp, error)
	GetAllUserID(context.Context, *GetAllUserIDReq) (*GetAllUserIDResp, error)
	AccountCheck(context.Context, *AccountCheckReq) (*AccountCheckResp, error)
	GetConversation(context.Context, *GetConversationReq) (*GetConversationResp, error)
	GetAllConversations(context.Context, *GetAllConversationsReq) (*GetAllConversationsResp, error)
	GetConversations(context.Context, *GetConversationsReq) (*GetConversationsResp, error)
	BatchSetConversations(context.Context, *BatchSetConversationsReq) (*BatchSetConversationsResp, error)
	SetConversation(context.Context, *SetConversationReq) (*SetConversationResp, error)
	SetRecvMsgOpt(context.Context, *SetRecvMsgOptReq) (*SetRecvMsgOptResp, error)
	GetUsers(context.Context, *GetUsersReq) (*GetUsersResp, error)
	AddUser(context.Context, *AddUserReq) (*AddUserResp, error)
	BlockUser(context.Context, *BlockUserReq) (*BlockUserResp, error)
	UnBlockUser(context.Context, *UnBlockUserReq) (*UnBlockUserResp, error)
	GetBlockUsers(context.Context, *GetBlockUsersReq) (*GetBlockUsersResp, error)
	AttributeSwitch(context.Context, *AttributeSwitchReq) (*AttributeSwitchResp, error)
	AttributeSwitchSet(context.Context, *AttributeSwitchSetReq) (*AttributeSwitchSetResp, error)
	AttributeMenu(context.Context, *AttributeMenuReq) (*AttributeMenuResp, error)
	// 用户反馈
	Feedback(context.Context, *FeedbackReq) (*FeedbackResp, error)
	// 常见问题
	GetCommonProblem(context.Context, *GetCommonProblemReq) (*GetCommonProblemResp, error)
	// 常见问题-反馈
	FeedbackCommonProblem(context.Context, *FeedbackCommonProblemReq) (*FeedbackCommonProblemResp, error)
	mustEmbedUnimplementedUserServer()
}

// UnimplementedUserServer must be embedded to have forward compatible implementations.
type UnimplementedUserServer struct {
}

func (UnimplementedUserServer) GetUserInfo(context.Context, *GetUserInfoReq) (*GetUserInfoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfo not implemented")
}
func (UnimplementedUserServer) UpdateUserInfo(context.Context, *UpdateUserInfoReq) (*UpdateUserInfoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserInfo not implemented")
}
func (UnimplementedUserServer) SetGlobalRecvMessageOpt(context.Context, *SetGlobalRecvMessageOptReq) (*SetGlobalRecvMessageOptResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetGlobalRecvMessageOpt not implemented")
}
func (UnimplementedUserServer) GetAllUserID(context.Context, *GetAllUserIDReq) (*GetAllUserIDResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllUserID not implemented")
}
func (UnimplementedUserServer) AccountCheck(context.Context, *AccountCheckReq) (*AccountCheckResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AccountCheck not implemented")
}
func (UnimplementedUserServer) GetConversation(context.Context, *GetConversationReq) (*GetConversationResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConversation not implemented")
}
func (UnimplementedUserServer) GetAllConversations(context.Context, *GetAllConversationsReq) (*GetAllConversationsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllConversations not implemented")
}
func (UnimplementedUserServer) GetConversations(context.Context, *GetConversationsReq) (*GetConversationsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConversations not implemented")
}
func (UnimplementedUserServer) BatchSetConversations(context.Context, *BatchSetConversationsReq) (*BatchSetConversationsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchSetConversations not implemented")
}
func (UnimplementedUserServer) SetConversation(context.Context, *SetConversationReq) (*SetConversationResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetConversation not implemented")
}
func (UnimplementedUserServer) SetRecvMsgOpt(context.Context, *SetRecvMsgOptReq) (*SetRecvMsgOptResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetRecvMsgOpt not implemented")
}
func (UnimplementedUserServer) GetUsers(context.Context, *GetUsersReq) (*GetUsersResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsers not implemented")
}
func (UnimplementedUserServer) AddUser(context.Context, *AddUserReq) (*AddUserResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUser not implemented")
}
func (UnimplementedUserServer) BlockUser(context.Context, *BlockUserReq) (*BlockUserResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BlockUser not implemented")
}
func (UnimplementedUserServer) UnBlockUser(context.Context, *UnBlockUserReq) (*UnBlockUserResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnBlockUser not implemented")
}
func (UnimplementedUserServer) GetBlockUsers(context.Context, *GetBlockUsersReq) (*GetBlockUsersResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockUsers not implemented")
}
func (UnimplementedUserServer) AttributeSwitch(context.Context, *AttributeSwitchReq) (*AttributeSwitchResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AttributeSwitch not implemented")
}
func (UnimplementedUserServer) AttributeSwitchSet(context.Context, *AttributeSwitchSetReq) (*AttributeSwitchSetResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AttributeSwitchSet not implemented")
}
func (UnimplementedUserServer) AttributeMenu(context.Context, *AttributeMenuReq) (*AttributeMenuResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AttributeMenu not implemented")
}
func (UnimplementedUserServer) Feedback(context.Context, *FeedbackReq) (*FeedbackResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Feedback not implemented")
}
func (UnimplementedUserServer) GetCommonProblem(context.Context, *GetCommonProblemReq) (*GetCommonProblemResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCommonProblem not implemented")
}
func (UnimplementedUserServer) FeedbackCommonProblem(context.Context, *FeedbackCommonProblemReq) (*FeedbackCommonProblemResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FeedbackCommonProblem not implemented")
}
func (UnimplementedUserServer) mustEmbedUnimplementedUserServer() {}

// UnsafeUserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServer will
// result in compilation errors.
type UnsafeUserServer interface {
	mustEmbedUnimplementedUserServer()
}

func RegisterUserServer(s grpc.ServiceRegistrar, srv UserServer) {
	s.RegisterService(&User_ServiceDesc, srv)
}

func _User_GetUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_GetUserInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetUserInfo(ctx, req.(*GetUserInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_UpdateUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).UpdateUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_UpdateUserInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).UpdateUserInfo(ctx, req.(*UpdateUserInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_SetGlobalRecvMessageOpt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetGlobalRecvMessageOptReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).SetGlobalRecvMessageOpt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_SetGlobalRecvMessageOpt_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).SetGlobalRecvMessageOpt(ctx, req.(*SetGlobalRecvMessageOptReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetAllUserID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllUserIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetAllUserID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_GetAllUserID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetAllUserID(ctx, req.(*GetAllUserIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_AccountCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccountCheckReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).AccountCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_AccountCheck_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).AccountCheck(ctx, req.(*AccountCheckReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetConversation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConversationReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetConversation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_GetConversation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetConversation(ctx, req.(*GetConversationReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetAllConversations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllConversationsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetAllConversations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_GetAllConversations_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetAllConversations(ctx, req.(*GetAllConversationsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetConversations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConversationsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetConversations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_GetConversations_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetConversations(ctx, req.(*GetConversationsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_BatchSetConversations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchSetConversationsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).BatchSetConversations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_BatchSetConversations_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).BatchSetConversations(ctx, req.(*BatchSetConversationsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_SetConversation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetConversationReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).SetConversation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_SetConversation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).SetConversation(ctx, req.(*SetConversationReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_SetRecvMsgOpt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRecvMsgOptReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).SetRecvMsgOpt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_SetRecvMsgOpt_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).SetRecvMsgOpt(ctx, req.(*SetRecvMsgOptReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsersReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_GetUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetUsers(ctx, req.(*GetUsersReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_AddUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddUserReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).AddUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_AddUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).AddUser(ctx, req.(*AddUserReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_BlockUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlockUserReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).BlockUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_BlockUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).BlockUser(ctx, req.(*BlockUserReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_UnBlockUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnBlockUserReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).UnBlockUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_UnBlockUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).UnBlockUser(ctx, req.(*UnBlockUserReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetBlockUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBlockUsersReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetBlockUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_GetBlockUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetBlockUsers(ctx, req.(*GetBlockUsersReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_AttributeSwitch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AttributeSwitchReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).AttributeSwitch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_AttributeSwitch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).AttributeSwitch(ctx, req.(*AttributeSwitchReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_AttributeSwitchSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AttributeSwitchSetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).AttributeSwitchSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_AttributeSwitchSet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).AttributeSwitchSet(ctx, req.(*AttributeSwitchSetReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_AttributeMenu_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AttributeMenuReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).AttributeMenu(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_AttributeMenu_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).AttributeMenu(ctx, req.(*AttributeMenuReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_Feedback_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FeedbackReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).Feedback(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_Feedback_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).Feedback(ctx, req.(*FeedbackReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetCommonProblem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCommonProblemReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetCommonProblem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_GetCommonProblem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetCommonProblem(ctx, req.(*GetCommonProblemReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_FeedbackCommonProblem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FeedbackCommonProblemReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).FeedbackCommonProblem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_FeedbackCommonProblem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).FeedbackCommonProblem(ctx, req.(*FeedbackCommonProblemReq))
	}
	return interceptor(ctx, in, info, handler)
}

// User_ServiceDesc is the grpc.ServiceDesc for User service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var User_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.user",
	HandlerType: (*UserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserInfo",
			Handler:    _User_GetUserInfo_Handler,
		},
		{
			MethodName: "UpdateUserInfo",
			Handler:    _User_UpdateUserInfo_Handler,
		},
		{
			MethodName: "SetGlobalRecvMessageOpt",
			Handler:    _User_SetGlobalRecvMessageOpt_Handler,
		},
		{
			MethodName: "GetAllUserID",
			Handler:    _User_GetAllUserID_Handler,
		},
		{
			MethodName: "AccountCheck",
			Handler:    _User_AccountCheck_Handler,
		},
		{
			MethodName: "GetConversation",
			Handler:    _User_GetConversation_Handler,
		},
		{
			MethodName: "GetAllConversations",
			Handler:    _User_GetAllConversations_Handler,
		},
		{
			MethodName: "GetConversations",
			Handler:    _User_GetConversations_Handler,
		},
		{
			MethodName: "BatchSetConversations",
			Handler:    _User_BatchSetConversations_Handler,
		},
		{
			MethodName: "SetConversation",
			Handler:    _User_SetConversation_Handler,
		},
		{
			MethodName: "SetRecvMsgOpt",
			Handler:    _User_SetRecvMsgOpt_Handler,
		},
		{
			MethodName: "GetUsers",
			Handler:    _User_GetUsers_Handler,
		},
		{
			MethodName: "AddUser",
			Handler:    _User_AddUser_Handler,
		},
		{
			MethodName: "BlockUser",
			Handler:    _User_BlockUser_Handler,
		},
		{
			MethodName: "UnBlockUser",
			Handler:    _User_UnBlockUser_Handler,
		},
		{
			MethodName: "GetBlockUsers",
			Handler:    _User_GetBlockUsers_Handler,
		},
		{
			MethodName: "AttributeSwitch",
			Handler:    _User_AttributeSwitch_Handler,
		},
		{
			MethodName: "AttributeSwitchSet",
			Handler:    _User_AttributeSwitchSet_Handler,
		},
		{
			MethodName: "AttributeMenu",
			Handler:    _User_AttributeMenu_Handler,
		},
		{
			MethodName: "Feedback",
			Handler:    _User_Feedback_Handler,
		},
		{
			MethodName: "GetCommonProblem",
			Handler:    _User_GetCommonProblem_Handler,
		},
		{
			MethodName: "FeedbackCommonProblem",
			Handler:    _User_FeedbackCommonProblem_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user/user.proto",
}