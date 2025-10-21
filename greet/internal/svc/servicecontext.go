package svc

import (
	"greet/internal/config"

	userpb "greet/proto/user"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type ServiceContext struct {
	Config  config.Config
	RPCuser userpb.UserClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	maxSize := c.MaxSize
	dialOption := grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxSize), grpc.MaxCallSendMsgSize(maxSize))
	opt := zrpc.WithDialOption(dialOption)
	return &ServiceContext{
		Config:  c,
		RPCuser: userpb.NewUserClient(zrpc.MustNewClient(c.RPCuser, opt).Conn()),
	}
}
