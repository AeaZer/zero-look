package errorx

import (
	"net/http"
)

// InvalidParam 参数错误
func InvalidParam() *BizError {
	return &BizError{
		code: http.StatusBadRequest,
		msg:  "参数错误",
	}
}
