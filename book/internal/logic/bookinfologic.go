// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"book/internal/svc"
	"book/internal/types"

	"github.com/bxcodec/faker/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type BookInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBookInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BookInfoLogic {
	return &BookInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type GreetResponse struct {
	Message string `json:"message"`
}

func (l *BookInfoLogic) BookInfo(req *types.BookInfoReq) (resp *types.BookInfoResp, err error) {
	// 调用 greet-api 服务获取问候语
	greetMsg := ""
	greetURL := fmt.Sprintf("%s/you", l.svcCtx.Config.GreetAPI.Endpoint)

	httpReq, err := http.NewRequestWithContext(l.ctx, http.MethodGet, greetURL, nil)
	if err != nil {
		logx.Errorf("failed to create request to greet-api: %v", err)
	} else {
		httpResp, err := l.svcCtx.GreetClient.Do(httpReq)
		if err != nil {
			logx.Errorf("failed to call greet-api: %v", err)
		} else {
			defer httpResp.Body.Close()
			body, err := io.ReadAll(httpResp.Body)
			if err != nil {
				logx.Errorf("failed to read greet-api response: %v", err)
			} else if httpResp.StatusCode == http.StatusOK {
				var greetResp GreetResponse
				if err := json.Unmarshal(body, &greetResp); err != nil {
					logx.Errorf("failed to unmarshal greet-api response: %v", err)
				} else {
					greetMsg = greetResp.Message
				}
			}
		}
	}

	resp = &types.BookInfoResp{
		Author: fmt.Sprintf("《%s》's author is %s. %s", req.Title, faker.Name(), greetMsg),
	}

	return
}
