package xerr

import "fmt"

type CodeError struct {
	errCode int
	errMsg string
}

func (e *CodeError) GetErrCode() int {
	return e.errCode
}

func (e *CodeError) GetErrMsg() string {
	return e.errMsg
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode: %d, ErrMsg: %s", e.errCode, e.errMsg)
}

func New(errCode int, errMsg string) *CodeError {
	return &CodeError{errCode, errMsg}
}

func NewErrCode(errCode int) *CodeError {
	return &CodeError{errCode: errCode, errMsg: MapErrMsg(errCode)}
}

func NewErrMsg(errMsg string) *CodeError {
	return &CodeError{errCode: BAD_REUQEST_ERROR, errMsg: errMsg}
}