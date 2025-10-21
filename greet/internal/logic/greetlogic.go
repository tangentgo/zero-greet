package logic

import (
	"context"
	"fmt"

	"greet/internal/svc"
	"greet/internal/types"
	"greet/proto/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GreetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGreetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GreetLogic {
	return &GreetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GreetLogic) Greet(req *types.Request) (resp *types.Response, err error) {
	resp = &types.Response{}
	rpcResp, err := l.svcCtx.RPCuser.Ping(l.ctx, &user.Request{
		Ping: req.Name,
	})
	resp.Message = fmt.Sprint(rpcResp.Pong + req.Name + " from greet-api")
	if err != nil {
		return
	}
	return
}
