package httpResult

import (
	"fmt"
	"net/http"

	"github.com/lerity-yao/server-result/xerr"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

func GetHttpErrCodeMsg(err error) (uint32, string, xerr.AlertLevel) {
	errCode := xerr.ServerCommonError
	errMsg := "服务器开小差啦，稍后再来试一试"
	var alertLevel xerr.AlertLevel

	// 追溯错误链中最初始的错误 （此处可以追溯出rpc服务的错误）
	// 所有逻辑中，最终返回的错误，都应该使用 errors.Wrapf()来返回错误
	causeErr := errors.Cause(err)
	// 断言err类型，看看是不是自己定义的err类型，如果是，就直接取自己定义的code和msg
	if e, ok := causeErr.(*xerr.CodeError); ok { //自定义错误类型
		//自定义CodeError
		errCode = e.GetErrCode()
		errMsg = e.GetErrMsg()
		alertLevel = e.GetAlertLevel()
	} else {
		// 只处理 grpc 的状态码和消息。通畅情况下，不会存在这种情况，因为 rpcResult.LoggerInterceptor的方法会对grpc的结果进行转换
		// 此处注意， rpc返回错误一定要用errors.Wrapf()，不然此处无法处理
		if grpcStatus, ok := status.FromError(causeErr); ok { // grpc err错误
			grpcMsg := grpcStatus.Message()
			grpcCode := grpcStatus.Code()

			if grpcMsg != "" {
				errMsg = grpcMsg
			}

			if grpcCode != 0 {
				errCode = uint32(grpcCode)
			}

		}
	}

	return errCode, errMsg, alertLevel
}

// buildErrFields 从 error 链构建日志字段（HttpResult / HttpStatusResult 用）
func buildErrFields(err error, alertLevel xerr.AlertLevel, result map[string]any) []logc.LogField {
	fields := []logc.LogField{
		logc.Field(xerr.LogFType, xerr.LogApiError),
		logc.Field(xerr.LogFResult, result),
		logc.Field(xerr.LogFStack, fmt.Sprintf("%+v", err)),
	}
	if alertLevel != "" {
		fields = append(fields, logc.Field(xerr.LogFAlertLevel, alertLevel))
	}
	return fields
}

// buildOptFields 从 opts 构建日志字段（ParamErrorResult / MdErrorResult / MapErrorResult 用）
func buildOptFields(result map[string]any, opts ...xerr.Option) []logc.LogField {
	fields := []logc.LogField{
		logc.Field(xerr.LogFType, xerr.LogApiError),
		logc.Field(xerr.LogFResult, result),
	}
	if al := xerr.ExtractAlertLevel(opts...); al != "" {
		fields = append(fields, logc.Field(xerr.LogFAlertLevel, al))
	}
	return fields
}

// HttpResult http返回结果
func HttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {

	// 获取当前链路跟踪trace和span
	traceId := trace.TraceIDFromContext(r.Context())
	spanId := trace.SpanIDFromContext(r.Context())
	// 返回成功
	if err == nil {

		logc.Infow(r.Context(), "ok",
			logc.Field(xerr.LogFType, xerr.LogApiSuccess),
			logc.Field(xerr.LogFResult, resp))
		r := Success(traceId, spanId, resp)
		httpx.WriteJson(w, http.StatusOK, r)
		return
	}

	//错误返回
	errCode, errMsg, alertLevel := GetHttpErrCodeMsg(err)

	logc.Errorw(r.Context(), errMsg, buildErrFields(err, alertLevel, map[string]any{"code": errCode, "msg": errMsg})...)
	httpx.WriteJson(w, http.StatusOK, Error(traceId, spanId, errCode, errMsg))
	return
}

// ParamErrorResult http 参数错误返回
func ParamErrorResult(r *http.Request, w http.ResponseWriter, err error, opts ...xerr.Option) {

	// 获取当前链路跟踪trace和span
	traceId := trace.TraceIDFromContext(r.Context())
	spanId := trace.SpanIDFromContext(r.Context())
	errMsg := xerr.MapErrMsg(xerr.RequestParamError) + " ," + err.Error()

	// 打印处理之后的错误
	logc.Errorw(r.Context(), errMsg, buildOptFields(map[string]any{"code": xerr.RequestParamError, "msg": errMsg}, opts...)...)
	httpx.WriteJson(w, http.StatusOK, Error(traceId, spanId, xerr.RequestParamError, errMsg))
}

// MdErrorResult 定义中间件错误
func MdErrorResult(r *http.Request, w http.ResponseWriter, msg string, opts ...xerr.Option) {
	// 获取当前链路跟踪trace和span
	traceId := trace.TraceIDFromContext(r.Context())
	spanId := trace.SpanIDFromContext(r.Context())
	errMsg := xerr.MapErrMsg(xerr.MdCommonError) + " ," + msg

	// 打印处理之后的错误
	logc.Errorw(r.Context(), errMsg, buildOptFields(map[string]any{"code": xerr.MdCommonError, "msg": errMsg}, opts...)...)
	httpx.WriteJson(w, http.StatusOK, Error(traceId, spanId, xerr.MdCommonError, errMsg))
}

// MapErrorResult 自定义返回code和msg
func MapErrorResult(r *http.Request, w http.ResponseWriter, code uint32, msg string, opts ...xerr.Option) {
	// 获取当前链路跟踪trace和span
	traceId := trace.TraceIDFromContext(r.Context())
	spanId := trace.SpanIDFromContext(r.Context())
	errMsg := xerr.MapErrMsg(code) + " ," + msg

	// 打印处理之后的错误
	logc.Errorw(r.Context(), errMsg, buildOptFields(map[string]any{"code": code, "msg": errMsg}, opts...)...)
	httpx.WriteJson(w, http.StatusOK, Error(traceId, spanId, code, errMsg))
}

// HttpStatusResult 返回自定义httpStatus状态码错误
func HttpStatusResult(r *http.Request, w http.ResponseWriter, statusCode int, err error) {

	// 获取当前链路跟踪trace和span
	traceId := trace.TraceIDFromContext(r.Context())
	spanId := trace.SpanIDFromContext(r.Context())

	//错误返回
	errCode, errMsg, alertLevel := GetHttpErrCodeMsg(err)

	logc.Errorw(r.Context(), errMsg, buildErrFields(err, alertLevel, map[string]any{"code": errCode, "msg": errMsg, "status": statusCode})...)

	httpx.WriteJson(w, statusCode, Error(traceId, spanId, errCode, errMsg))
}
