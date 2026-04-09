# server-result

go-zero 微服务统一响应与错误处理组件，提供 HTTP 响应封装、RPC 日志拦截、错误码管理和分级告警能力。

## 安装

```bash
go get github.com/lerity-yao/server-result
```

## 模块概览

```
server-result/
├── xerr/                  # 错误码定义 & CodeError 错误体系
│   ├── errMsg.go          # 错误码常量、告警级别常量、日志字段常量
│   └── erros.go           # CodeError 结构体、Option 模式、构造函数
├── httpResult/            # HTTP 统一响应
│   ├── httpResult.go      # 5 种响应方法 + 结构化日志
│   └── responseType.go    # 响应体结构定义
└── rpcResult/             # RPC 日志拦截
    └── loggerInterceptor.go  # gRPC unary 拦截器，结构化日志 + 错误转换
```

## 快速开始

### HTTP 响应

```go
// handler 中使用
func GetUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        resp, err := logic.NewGetUserLogic(r.Context(), svcCtx).GetUser(&req)
        httpResult.HttpResult(r, w, resp, err)
    }
}
```

### RPC 拦截器

```go
// main.go 中注册
s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
    // 注册服务...
})
s.AddUnaryInterceptors(rpcResult.LoggerInterceptor)
```

### 错误构造

```go
// 自定义 code + msg
xerr.NewErrCodeMsg(100010, "用户不存在")

// 用预定义 code，自动查 msg
xerr.NewErrCode(xerr.DbError)

// 固定 ServerCommonError code，自定义 msg
xerr.NewErrMsg("操作失败")
```

## 告警级别

支持 P0~P3 四级告警，通过 `WithAlertLevel` Option 注入，日志输出 `xr_alert_level` 字段供 Loki 查询。

```go
// 构造 CodeError 时指定告警级别
xerr.NewErrCodeMsg(100010, "核心服务异常", xerr.WithAlertLevel(xerr.AlertP0))

// 显式调用方法中指定
httpResult.ParamErrorResult(r, w, err, xerr.WithAlertLevel(xerr.AlertP2))
httpResult.MdErrorResult(r, w, "token过期", xerr.WithAlertLevel(xerr.AlertP1))
httpResult.MapErrorResult(r, w, code, msg, xerr.WithAlertLevel(xerr.AlertP3))
```

| 级别 | 常量 | 场景 |
|------|------|------|
| P0 | `xerr.AlertP0` | 核心服务不可用、未知错误（RPC 非 CodeError 自动标记） |
| P1 | `xerr.AlertP1` | 重要功能异常 |
| P2 | `xerr.AlertP2` | 一般功能异常 |
| P3 | `xerr.AlertP3` | 低优先级问题 |

> `WithAlertLevel` 参数类型为 `AlertLevel`，只能传 `AlertP0`~`AlertP3` 常量。

## HTTP 响应方法

| 方法 | 用途 | alertLevel 来源 |
|------|------|-----------------|
| `HttpResult` | 通用响应（成功/失败） | 自动从 CodeError 提取 |
| `HttpStatusResult` | 自定义 HTTP 状态码 | 自动从 CodeError 提取 |
| `ParamErrorResult` | 参数校验错误 | `...xerr.Option` 显式传入 |
| `MdErrorResult` | 中间件错误 | `...xerr.Option` 显式传入 |
| `MapErrorResult` | 自定义 code + msg | `...xerr.Option` 显式传入 |

## 结构化日志字段

所有日志字段使用 `xr_` 前缀，JSON 格式输出，支持 Loki LogQL 精准查询。

| 字段 | 说明 | 示例值 |
|------|------|--------|
| `xr_type` | 结果类型 | `API-SUCCESS` / `API-ERROR` / `RPC-SUCCESS` / `RPC-ERROR` |
| `xr_result` | 响应/错误摘要 | `{"code":100010,"msg":"用户不存在"}` |
| `xr_stack` | 错误堆栈 | 完整 `%+v` 堆栈信息 |
| `xr_alert_level` | 告警级别 | `P0` / `P1` / `P2` / `P3`（无告警时不输出） |

### Loki 查询示例

```logql
# 查所有 P0 告警
{app="my-service"} | json | xr_alert_level="P0"

# 查 API 错误
{app="my-service"} | json | xr_type="API-ERROR"

# 查 RPC 错误中的 P1 告警
{app="my-service"} | json | xr_type="RPC-ERROR" | xr_alert_level="P1"
```

## 预定义错误码

| 常量 | 值 | 说明 |
|------|----|------|
| `OK` | 0 | 成功 |
| `ServerCommonError` | 100001 | 服务器通用错误 |
| `RequestParamError` | 100002 | 参数错误 |
| `TokenExpireError` | 100003 | token 失效 |
| `TokenGenerateError` | 100004 | token 生成失败 |
| `DbError` | 100005 | 数据库错误 |
| `DbUpdateAffectedZeroError` | 100006 | 更新影响行数为 0 |
| `MdCommonError` | 100007 | 中间件错误 |
| `PermitNoAccess` | 100008 | 无权限 |
| `SignParamError` | 100009 | 签名错误 |
