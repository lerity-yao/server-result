# Changelog

## v0.0.6

### Features

- **WithAlertData 告警数据注入**：新增 `xerr.WithAlertData(key, value)` Option，支持逐个传入键值对，多次调用累加字段，日志统一输出 `xr_alert_data` 字段供 Loki 提取
- **AlertDataField 便捷函数**：新增 `xerr.AlertDataField(data)` 函数，返回 `logc.LogField`，支持在 `logc.Infow` / `logc.Errorw` 中直接传入告警数据
- **ExtractAlertData 工具函数**：新增 `xerr.ExtractAlertData(opts ...Option)` 函数，从 Option 中提取 alertData，供 httpResult 显式调用方法使用
- **HTTP 响应告警数据支持**：`HttpResult`、`HttpStatusResult` 自动从 CodeError 提取告警数据；`ParamErrorResult`、`MdErrorResult`、`MapErrorResult` 通过 `...xerr.Option` 显式传入
- **RPC 拦截器告警数据支持**：`LoggerInterceptor` 自动从 CodeError 提取并输出 `xr_alert_data` 字段
- **单元测试**：新增 httpResult 3 项 WithAlertData 测试用例

## v0.0.5

### Features

- **AlertField 便捷函数**：新增 `xerr.AlertField(level)` 函数，返回 `logc.LogField`，支持在 `logc.Infow` / `logc.Errorw` 中直接传入告警级别，无需手动拼装 `logc.Field(xerr.LogFAlertLevel, ...)`

## v0.0.4

### Features

- **告警级别体系**：新增 `AlertLevel` 自定义类型与 `AlertP0`~`AlertP3` 常量，`WithAlertLevel` 参数类型安全，防止传入任意字符串
- **Option 模式扩展 CodeError**：通过 `WithAlertLevel` Option 注入告警级别，构造函数向后兼容
- **结构化日志**：所有日志字段统一 `xr_` 前缀（`xr_type`、`xr_result`、`xr_stack`、`xr_alert_level`），JSON 格式输出，支持 Loki LogQL 精准查询
- **HTTP 响应告警支持**：`HttpResult`、`HttpStatusResult` 自动从 CodeError 提取告警级别；`ParamErrorResult`、`MdErrorResult`、`MapErrorResult` 支持 `...xerr.Option` 显式传入
- **RPC 拦截器改造**：`LoggerInterceptor` 输出结构化字段，成功请求记录 `RPC-SUCCESS`，错误请求记录完整错误摘要与堆栈
- **未知错误强制 P0 告警**：RPC 拦截器中非 `CodeError` 的错误自动标记 `xr_alert_level: P0`
- **单元测试**：新增 httpResult 12 项测试用例，覆盖全部响应方法及告警级别场景

### Docs

- 新增 README.md 项目文档

## v0.0.3

- fix: 修复 grpc code 漏掉逻辑 bug

## v0.0.2

- fix: 调整 go version

## v0.0.1

- feat: 项目初始化，基础错误码与 HTTP 响应封装
