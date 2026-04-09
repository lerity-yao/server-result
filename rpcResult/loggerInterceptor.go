package rpcResult

import (
	"context"
	"fmt"

	"github.com/lerity-yao/server-result/xerr"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggerInterceptor 拦截处理 zrpc 日志， 作为 given unary interceptors 加入 zrpc服务
// 把 zrpc错误转成自定义错误,所以返回错误的时候需要使用errors.Wrapf()来返回错误，
// 如果没有这么写，后续没有转换，会引起连锁反应，导致 HttpResult 无法正确捕获错
func LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err == nil {
		logc.Infow(ctx, "ok",
			logc.Field(xerr.LogFType, xerr.LogRpcSuccess),
			logc.Field(xerr.LogFResult, resp))
		return resp, nil
	}

	// 追溯错误链中最初始的错误
	causeErr := errors.Cause(err)
	// 断言err类型是否自定义类型
	if e, ok := causeErr.(*xerr.CodeError); ok { //自定义错误类型
		fields := []logc.LogField{
			logc.Field(xerr.LogFType, xerr.LogRpcError),
			logc.Field(xerr.LogFResult, map[string]any{"code": e.GetErrCode(), "msg": e.GetErrMsg()}),
			logc.Field(xerr.LogFStack, fmt.Sprintf("%+v", err)),
		}
		if al := e.GetAlertLevel(); al != "" {
			fields = append(fields, logc.Field(xerr.LogFAlertLevel, al))
		}
		logc.Errorw(ctx, e.GetErrMsg(), fields...)

		//把 zrpc错误转成自定义错误,所以返回错误的时候需要使用errors.Wrapf()来返回错误
		err = status.Error(codes.Code(e.GetErrCode()), e.GetErrMsg())
	} else {
		logc.Errorw(ctx, causeErr.Error(),
			logc.Field(xerr.LogFType, xerr.LogRpcError),
			logc.Field(xerr.LogFResult, map[string]any{"msg": causeErr.Error()}),
			logc.Field(xerr.LogFStack, fmt.Sprintf("%+v", err)),
			logc.Field(xerr.LogFAlertLevel, xerr.AlertP0))
	}

	return resp, err
}
