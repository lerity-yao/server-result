package rpcResult

import (
	"context"
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
	if err != nil {
		// 追溯错误链中最初始的错误
		causeErr := errors.Cause(err)
		// 断言err类型是否自定义类型
		if e, ok := causeErr.(*xerr.CodeError); ok { //自定义错误类型

			logc.Errorf(ctx, "[RPC-ERR]: %+v", err)

			//把 zrpc错误转成自定义错误,所以返回错误的时候需要使用errors.Wrapf()来返回错误
			err = status.Error(codes.Code(e.GetErrCode()), e.GetErrMsg())
		} else {
			logc.Errorf(ctx, "[RPC-ERR]: %+v", err)
		}

	}
	return resp, err
}
