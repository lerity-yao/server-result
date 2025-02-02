package xerr

// OK 成功返回
const OK uint32 = 0

// ServerCommonError 服务器开小差啦,稍后再来试一试
const ServerCommonError uint32 = 100001

// RequestParamError 参数错误
const RequestParamError uint32 = 100002

// TokenExpireError token失效，请重新登陆
const TokenExpireError uint32 = 100003

// TokenGenerateError 生成token失败
const TokenGenerateError uint32 = 100004

// DbError 数据库繁忙,请稍后再试
const DbError uint32 = 100005

// DbUpdateAffectedZeroError 更新数据影响行数为0
const DbUpdateAffectedZeroError uint32 = 100006

// MdCommonError 中间件错误
const MdCommonError uint32 = 100007

// PermitNoAccess 无权限操作
const PermitNoAccess uint32 = 100008

// SignParamError 签名错误
const SignParamError uint32 = 100009

var message map[uint32]string

// init 初始化 message map
func init() {
	message = make(map[uint32]string)
	message[OK] = "SUCCESS"
	message[ServerCommonError] = "服务器开小差啦,稍后再来试一试"
	message[RequestParamError] = "参数错误"
	message[TokenExpireError] = "token失效，请重新登陆"
	message[TokenGenerateError] = "生成token失败"
	message[DbError] = "数据库繁忙,请稍后再试"
	message[DbUpdateAffectedZeroError] = "更新数据影响行数为0"
	message[MdCommonError] = "中间件错误"
	message[PermitNoAccess] = "无权限操作"
	message[SignParamError] = "签名错误"
}

// MapErrMsg 查询errCode对应的msg, errCode 只能是errCode中定义的，不然返回 服务器开小差啦,稍后再来试一试
func MapErrMsg(errCode uint32) string {
	if msg, ok := message[errCode]; ok {
		return msg
	} else {
		return "服务器开小差啦,稍后再来试一试"
	}
}

// IsCodeErr 判断是否全局错误
func IsCodeErr(errCode uint32) bool {
	if _, ok := message[errCode]; ok {
		return true
	} else {
		return false
	}
}
