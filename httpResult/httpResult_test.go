package httpResult

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lerity-yao/server-result/xerr"
	"github.com/pkg/errors"
)

// 模拟业务 resp
type UserResp struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestHttpResult_Success(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/user/1", nil)
	r = r.WithContext(context.Background())
	w := httptest.NewRecorder()

	resp := &UserResp{Id: 1001, Name: "张三", Age: 25}

	HttpResult(r, w, resp, nil)

	t.Logf("HTTP Status: %d", w.Code)
	t.Logf("Response Body: %s", w.Body.String())
}

func TestHttpResult_CodeError(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/user/create", nil)
	r = r.WithContext(context.Background())
	w := httptest.NewRecorder()

	err := errors.Wrapf(xerr.NewErrCodeMsg(100010, "用户不存在"), "查询用户失败 uid: %d", 1001)
	HttpResult(r, w, nil, err)

	t.Logf("HTTP Status: %d", w.Code)
	t.Logf("Response Body: %s", w.Body.String())
}

func TestHttpResult_CodeError_WithAlert(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/user/create", nil)
	r = r.WithContext(context.Background())
	w := httptest.NewRecorder()

	err := errors.Wrapf(xerr.NewErrCodeMsg(100010, "用户不存在", xerr.WithAlertLevel(xerr.AlertP1)), "查询用户失败 uid: %d", 1001)
	HttpResult(r, w, nil, err)

	t.Logf("HTTP Status: %d", w.Code)
	t.Logf("Response Body: %s", w.Body.String())
}

func TestHttpResult_ServerError(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/order/create", nil)
	r = r.WithContext(context.Background())
	w := httptest.NewRecorder()

	err := errors.Wrap(errors.New("connection refused"), "调用订单RPC失败")
	HttpResult(r, w, nil, err)

	t.Logf("HTTP Status: %d", w.Code)
	t.Logf("Response Body: %s", w.Body.String())
}

func TestParamErrorResult(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/user/create", nil)
	r = r.WithContext(context.Background())
	w := httptest.NewRecorder()

	err := errors.New("name is required")
	ParamErrorResult(r, w, err)

	t.Logf("HTTP Status: %d", w.Code)
	t.Logf("Response Body: %s", w.Body.String())
}

func TestParamErrorResult_WithAlert(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/user/create", nil)
	r = r.WithContext(context.Background())
	w := httptest.NewRecorder()

	err := errors.New("name is required")
	ParamErrorResult(r, w, err, xerr.WithAlertLevel(xerr.AlertP2))

	t.Logf("HTTP Status: %d", w.Code)
	t.Logf("Response Body: %s", w.Body.String())
}

func TestMdErrorResult(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/user/1", nil)
	r = r.WithContext(context.Background())
	w := httptest.NewRecorder()

	MdErrorResult(r, w, "token已过期")

	t.Logf("HTTP Status: %d", w.Code)
	t.Logf("Response Body: %s", w.Body.String())
}

func TestMdErrorResult_WithAlert(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/user/1", nil)
	r = r.WithContext(context.Background())
	w := httptest.NewRecorder()

	MdErrorResult(r, w, "token已过期", xerr.WithAlertLevel(xerr.AlertP2))

	t.Logf("HTTP Status: %d", w.Code)
	t.Logf("Response Body: %s", w.Body.String())
}

func TestMapErrorResult(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/order/pay", nil)
	r = r.WithContext(context.Background())
	w := httptest.NewRecorder()

	MapErrorResult(r, w, xerr.PermitNoAccess, "无法操作他人订单")

	t.Logf("HTTP Status: %d", w.Code)
	t.Logf("Response Body: %s", w.Body.String())
}

func TestMapErrorResult_WithAlert(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/order/pay", nil)
	r = r.WithContext(context.Background())
	w := httptest.NewRecorder()

	MapErrorResult(r, w, xerr.PermitNoAccess, "无法操作他人订单", xerr.WithAlertLevel(xerr.AlertP1))

	t.Logf("HTTP Status: %d", w.Code)
	t.Logf("Response Body: %s", w.Body.String())
}

func TestHttpStatusResult(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/admin/dashboard", nil)
	r = r.WithContext(context.Background())
	w := httptest.NewRecorder()

	err := errors.Wrapf(xerr.NewErrCodeMsg(xerr.TokenExpireError, "token失效，请重新登陆"), "鉴权失败")
	HttpStatusResult(r, w, http.StatusUnauthorized, err)

	t.Logf("HTTP Status: %d", w.Code)
	t.Logf("Response Body: %s", w.Body.String())
}

func TestHttpStatusResult_WithAlert(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/admin/dashboard", nil)
	r = r.WithContext(context.Background())
	w := httptest.NewRecorder()

	err := errors.Wrapf(xerr.NewErrCodeMsg(xerr.TokenExpireError, "token失效，请重新登陆", xerr.WithAlertLevel(xerr.AlertP0)), "鉴权失败")
	HttpStatusResult(r, w, http.StatusUnauthorized, err)

	t.Logf("HTTP Status: %d", w.Code)
	t.Logf("Response Body: %s", w.Body.String())
}
