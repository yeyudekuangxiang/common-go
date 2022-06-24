package errno

import (
	"fmt"
	"runtime"
	"strconv"
)

/*
状态码参考 https://segmentfault.com/a/1190000020525050
第一位表示错误级别, 1 为系统错误, 2 为普通错误
第二三位表示服务模块代码
第四五位表示具体错误代码
*/
var (
	OK = Err{code: 200, message: "OK"}

	// ErrInternalServer 系统错误
	ErrInternalServer = Err{code: 10001, message: "内部服务器错误"}
	// ErrBind 绑定错误
	ErrBind = Err{code: 10002, message: "请求参数错误"}
	// ErrLimit 超出频率限制
	ErrLimit   = Err{code: 10003, message: "操作太频繁了、请稍后再试"}
	ErrTimeout = Err{code: 10004, message: "操作已超时"}

	// ErrRecordNotFound 数据库错误
	ErrRecordNotFound = Err{code: 20100, message: "数据异常"}

	// ErrAuth 未登录
	ErrAuth = Err{code: 20201, message: "未登陆"}
	// ErrValidation 验证失败
	ErrValidation = Err{code: 20202, message: "验证失败"}

	// ErrUserNotFound 未查询到用户信息
	ErrUserNotFound = Err{code: 20301, message: "未查询到用户信息"}

	ErrNotBindMobile = Err{code: 20303, message: "未授权手机号码"}

	// ErrAdminNotFound 管理员错误 前缀204
	ErrAdminNotFound = Err{code: 20401, message: "未查询到管理员信息"}
)

// Err 定义错误
type Err struct {
	code    int    // 错误码
	message string // 展示给用户看的
	err     error  // 保存内部错误信息
	callers string //保存调用文件名
}

func (e Err) WithErr(err error) Err {
	e.err = err
	return e
}
func (e Err) WithCaller() Err {
	_, f, l, _ := runtime.Caller(1)
	e.callers = f + ":" + strconv.Itoa(l)
	return e
}
func (e Err) With(err error) Err {
	e.err = err
	_, f, l, _ := runtime.Caller(1)
	e.callers = f + ":" + strconv.Itoa(l)
	return e
}
func (e Err) Code() int {
	return e.code
}
func (e Err) Message() string {
	return e.message
}
func (e Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s", e.code, e.message)
}

// DecodeErr 解码错误, 获取 Code 和 Message
func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code(), OK.Message()
	}
	if e, ok := err.(Err); ok {
		return e.Code(), e.Message()
	}
	return ErrInternalServer.Code(), err.Error()
}
