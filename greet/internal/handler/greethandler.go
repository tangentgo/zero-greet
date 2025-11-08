package handler

import (
	"context"
	"net/http"

	"greet/internal/logic"
	"greet/internal/svc"
	"greet/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GreetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 传播追踪头
		ctx := propagateTracingHeaders(r)

		l := logic.NewGreetLogic(ctx, svcCtx)
		resp, err := l.Greet(&req)
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
