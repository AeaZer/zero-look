package errorx

import (
	"net/http"
)

// InvalidParam 参数错误
func InvalidParam() *BizError {
	return New(http.StatusBadRequest, "参数错误")
}
