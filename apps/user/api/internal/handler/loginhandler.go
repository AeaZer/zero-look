package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/zero-look/apps/user/api/internal/logic"
	"github.com/zero-look/apps/user/api/internal/svc"
	"github.com/zero-look/apps/user/api/internal/types"
)

func loginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
