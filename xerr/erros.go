package xerr

import "strconv"

// CodeError 常用通用固定错误
type CodeError struct {
	errCode    uint32
	errMsg     string
	alertLevel AlertLevel
}

type Option func(*CodeError)

func WithAlertLevel(level AlertLevel) Option {
	return func(e *CodeError) { e.alertLevel = level }
}

func (e *CodeError) GetAlertLevel() AlertLevel { return e.alertLevel }

func applyOpts(e *CodeError, opts []Option) {
	for _, opt := range opts {
		opt(e)
	}
}

// ExtractAlertLevel 从 Option 中提取 alertLevel，供 httpResult 显式调用方法使用
func ExtractAlertLevel(opts ...Option) AlertLevel {
	if len(opts) == 0 {
		return ""
	}
	e := &CodeError{}
	applyOpts(e, opts)
	return e.alertLevel
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
	return "ErrCode: " + strconv.Itoa(int(e.errCode)) + ", ErrMsg: " + e.errMsg
}

// NewErrCodeMsg 构建新的ErrCodeMsg，允许自定义 errCode 和 errMsg
func NewErrCodeMsg(errCode uint32, errMsg string, opts ...Option) *CodeError {
	e := &CodeError{errCode: errCode, errMsg: errMsg}
	if len(opts) > 0 {
		applyOpts(e, opts)
	}
	return e
}

// NewErrCode 用自定义的 errCode 来 map message 寻找已有的errMsg，没有则返回 服务器开小差啦,稍后再来试一试
func NewErrCode(errCode uint32, opts ...Option) *CodeError {
	e := &CodeError{errCode: errCode, errMsg: MapErrMsg(errCode)}
	if len(opts) > 0 {
		applyOpts(e, opts)
	}
	return e
}

// NewErrMsg 固定errCode为 ServerCommonError,自定义errMsg
func NewErrMsg(errMsg string, opts ...Option) *CodeError {
	e := &CodeError{errCode: ServerCommonError, errMsg: errMsg}
	if len(opts) > 0 {
		applyOpts(e, opts)
	}
	return e
}
