package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/zero-look/apps/user/api/internal/logic"
	"github.com/zero-look/apps/user/api/internal/svc"
	"github.com/zero-look/apps/user/api/internal/types"
)

func registerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRegisterLogic(r.Context(), svcCtx)
		err := l.Register(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
