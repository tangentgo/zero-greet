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

// tracingHeadersKey 是用于在 context 中存储 HTTP headers 的 key
type tracingHeadersKey struct{}

func BookInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BookInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 将整个 HTTP request 存储到 context 中，以便 logic 层可以访问 headers
		ctx := context.WithValue(r.Context(), tracingHeadersKey{}, r.Header)

		l := logic.NewBookInfoLogic(ctx, svcCtx)
		resp, err := l.BookInfo(&req)
		if err != nil {
			httpx.ErrorCtx(ctx, w, err)
		} else {
			httpx.OkJsonCtx(ctx, w, resp)
		}
	}
}
