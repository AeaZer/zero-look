package errorx

type BizError struct {
	code int
	msg  string
}

func (b *BizError) Error() string {
	return b.msg
}

func New(code int, msg string) *BizError {
	return &BizError{
		code: code,
		msg:  msg,
	}
}

type errorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (b *BizError) Data() *errorResponse {
	return &errorResponse{
		Code: b.code,
		Msg:  b.msg,
	}
}

func (e *errorResponse) AppendMicroSerialNumber(bizNumber int) {
	e.Code += bizNumber * 1e4
}
