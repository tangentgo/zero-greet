package logic

import (
	"context"

	"greet/internal/svc"
	"greet/internal/types"
	"greet/proto/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type BaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseLogic {
	return &BaseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseLogic) Base(req *types.BaseHandlerReq) (resp *types.BaseHandlerResp, err error) {
	resp = &types.BaseHandlerResp{}
	rpcResp,err := l.svcCtx.RPCuser.Ping(l.ctx,&user.Request{
		Ping: "ping from greet",
	})

	if err != nil {
		return resp,err
	}

	if rpcResp == nil {
		return resp,nil
	}

	resp.Message = rpcResp.GetPong() + " from greet"

	logx.Infof("rpc resp: %v",rpcResp.GetPong())
	return
}
