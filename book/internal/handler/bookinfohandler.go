// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"context"
	"net/http"

	"book/internal/logic"
	"book/internal/svc"
	"book/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func BookInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BookInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 将追踪头存储到 context 中
		ctx := propagateTracingHeaders(r)

		l := logic.NewBookInfoLogic(ctx, svcCtx)
		resp, err := l.BookInfo(&req)
		if err != nil {
			httpx.ErrorCtx(ctx, w, err)
		} else {
			httpx.OkJsonCtx(ctx, w, resp)
		}
	}
}

// propagateTracingHeaders 从请求中提取追踪头并存储到 context
func propagateTracingHeaders(r *http.Request) context.Context {
	ctx := r.Context()

	// Istio/Zipkin B3 追踪头
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
		if value := r.Header.Get(header); value != "" {
			ctx = context.WithValue(ctx, header, value)
		}
	}

	return ctx
}
