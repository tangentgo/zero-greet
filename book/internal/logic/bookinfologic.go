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
		// 传播追踪头 - Istio 需要这些头来关联请求
		// go-zero 会自动从 context 中提取并设置这些头
		l.propagateTracingHeaders(httpReq)

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

// tracingHeadersKey 是用于在 context 中存储 HTTP headers 的 key
type tracingHeadersKey struct{}

// propagateTracingHeaders 从 context 中获取追踪头并设置到 HTTP 请求
func (l *BookInfoLogic) propagateTracingHeaders(req *http.Request) {
	// 从 context 中获取原始 HTTP headers
	if headers := l.ctx.Value(tracingHeadersKey{}); headers != nil {
		if httpHeaders, ok := headers.(http.Header); ok {
			tracingHeaders := []string{
				"x-request-id",
				"x-b3-traceid",
				"x-b3-spanid",
				"x-b3-parentspanid",
				"x-b3-sampled",
				"x-b3-flags",
				"b3",
				"x-ot-span-context",
				"x-cloud-trace-context",
				"traceparent",
				"grpc-trace-bin",
			}

			for _, header := range tracingHeaders {
				if value := httpHeaders.Get(header); value != "" {
					req.Header.Set(header, value)
				}
			}
		}
	}
}
