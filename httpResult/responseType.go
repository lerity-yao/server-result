package httpResult

import "github.com/lerity-yao/server-result/xerr"

type ResponseSuccessBean struct {
	Code    uint32      `json:"code"`    // 业务状态码
	Msg     string      `json:"msg"`     // 业务消息
	Data    interface{} `json:"data"`    // 返回数据
	TraceId string      `json:"traceId"` // 链路跟踪traceId
	SpanId  string      `json:"spanId"`  // 链路跟踪spanId
}
type NullJson struct{}

// Success 请求成功返回数据, traceId 为链路跟踪traceId, spanId为链路跟踪spanId
func Success(traceId, spanId string, data interface{}) *ResponseSuccessBean {
	return &ResponseSuccessBean{
		Code:    xerr.OK,
		Msg:     xerr.MapErrMsg(xerr.OK),
		Data:    data,
		TraceId: traceId,
		SpanId:  spanId,
	}
}

type ResponseErrorBean struct {
	Code    uint32 `json:"code"`
	Msg     string `json:"msg"`
	TraceId string `json:"traceId"` // 链路跟踪traceId
	SpanId  string `json:"spanId"`  // 链路跟踪spanId
}

// Error 请求失败返回数据, traceId 为链路跟踪traceId, spanId为链路跟踪spanId
func Error(traceId, spanId string, errCode uint32, errMsg string) *ResponseErrorBean {
	return &ResponseErrorBean{
		Code:    errCode,
		Msg:     errMsg,
		TraceId: traceId,
		SpanId:  spanId,
	}
}
