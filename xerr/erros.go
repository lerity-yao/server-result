package xerr

import "fmt"

// CodeError 常用通用固定错误
type CodeError struct {
	errCode uint32
	errMsg  string
}

// GetErrCode 返回给前端的错误码
func (e *CodeError) GetErrCode() uint32 {
	return e.errCode
}

// GetErrMsg 返回给前端显示端错误信息
func (e *CodeError) GetErrMsg() string {
	return e.errMsg
}

// Error 构建error接口，返回字符串，可以允许错误类型断言
func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode: %d, ErrMsg: %s", e.errCode, e.errMsg)
}

// NewErrCodeMsg 构建新的ErrCodeMsg，允许自定义 errCode 和 errMsg
func NewErrCodeMsg(errCode uint32, errMsg string) *CodeError {
	return &CodeError{errCode: errCode, errMsg: errMsg}
}

// NewErrCode 用自定义的 errCode 来 map message 寻找已有的errMsg，没有则返回 服务器开小差啦,稍后再来试一试
func NewErrCode(errCode uint32) *CodeError {
	return &CodeError{errCode: errCode, errMsg: MapErrMsg(errCode)}
}

// NewErrMsg 固定errCode为 ServerCommonError,自定义errMsg
func NewErrMsg(errMsg string) *CodeError {
	return &CodeError{errCode: ServerCommonError, errMsg: errMsg}
}
