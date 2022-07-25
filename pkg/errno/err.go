package errno

import (
	"errors"
	"fmt"
	"mio/config"
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
	OK = err{code: 200, message: "OK"}

	// ErrInternalServer 系统错误
	ErrInternalServer = err{code: 10001, message: "内部服务器错误"}
	// ErrBind 绑定错误
	ErrBind = err{code: 10002, message: "请求参数错误"}
	// ErrLimit 超出频率限制
	ErrLimit   = err{code: 10003, message: "操作太频繁了、请稍后再试"}
	ErrTimeout = err{code: 10004, message: "操作已超时"}

	// ErrRecordNotFound 数据库错误
	ErrRecordNotFound = err{code: 20100, message: "数据异常"}

	// ErrAuth 未登录
	ErrAuth = err{code: 20201, message: "未登陆"}
	// ErrValidation 验证失败
	ErrValidation = err{code: 20202, message: "验证失败"}

	// ErrUserNotFound 未查询到用户信息
	ErrUserNotFound = err{code: 20301, message: "未查询到用户信息"}

	// ErrNotBindMobile 未绑定手机号
	ErrNotBindMobile = err{code: 20303, message: "未授权手机号码"}
	// ErrBindMobile 绑定手机号时异常
	ErrBindMobile = err{code: 20304, message: "绑定手机号码失败"}

	// ErrAdminNotFound 管理员错误 前缀204
	ErrAdminNotFound = err{code: 20401, message: "未查询到管理员信息"}
)

// Err 定义错误
type err struct {
	code    int    // 错误码
	message string // 展示给用户看的
	err     error  // 保存内部错误信息
	callers string //保存调用文件名
}

// WithErr 带上err信息
func (e err) WithErr(err error) err {
	e.err = err
	return e
}

// WithMessage 替换默认的提示
func (e err) WithMessage(message string) err {
	e.message = message
	return e
}

// WithErrMessage 带上err message
func (e err) WithErrMessage(err string) err {
	e.err = errors.New(err)
	return e
}

// WithCaller 带上调用栈
func (e err) WithCaller() err {
	_, f, l, _ := runtime.Caller(1)
	e.callers = f + ":" + strconv.Itoa(l)
	return e
}

// With 带上错误和调用栈
func (e err) With(err error) err {
	e.err = err
	_, f, l, _ := runtime.Caller(1)
	e.callers = f + ":" + strconv.Itoa(l)
	return e
}
func (e err) Code() int {
	return e.code
}
func (e err) Message() string {
	return e.message
}
func (e err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s ,err: %v", e.code, e.message, e.err)
}

// DecodeErr 解码错误, 获取 Code 和 Message
func DecodeErr(e error) (int, string) {
	if e == nil {
		return OK.Code(), OK.Message()
	}
	if decodeErr, ok := e.(err); ok {
		if config.Config.App.Debug {
			return decodeErr.Code(), decodeErr.Error()
		}
		return decodeErr.Code(), decodeErr.Message()
	}
	if config.Config.App.Debug {
		return ErrInternalServer.Code(), e.Error()
	}
	return ErrInternalServer.Code(), ErrInternalServer.Error()
	//后面系统全面替换后使用下面的方式
	//return ErrInternalServer.Code(), ErrInternalServer.Message()
}
