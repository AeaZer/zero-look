package api

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/zero-look/pkg/errorx"
)

func SetGlobalErrorHandler(bizNumber int) {
	httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, any) {
		switch e := err.(type) {
		case *errorx.BizError:
			errorResponse := e.Data()
			errorResponse.AppendMicroSerialNumber(bizNumber)
			return http.StatusOK, errorResponse
		default:
			return http.StatusInternalServerError, nil
		}
	})
}
