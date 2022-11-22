package errno

import (
	"errors"
	"fmt"
	"log"
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

	// ErrCommon 通用错误
	ErrCommon = err{code: 10000, message: "系统繁忙,请稍后再试"}
	// ErrInternalServer 系统错误
	ErrInternalServer = err{code: 10001, message: "内部服务器错误"}
	// ErrValidation 参数校验错误
	ErrValidation = err{code: 10002, message: "请求参数错误"}
	// ErrLimit 超出频率限制
	ErrLimit   = err{code: 10003, message: "操作太频繁了、请稍后再试"}
	ErrTimeout = err{code: 10004, message: "操作已超时"}

	// ErrRecordNotFound 数据库错误
	ErrRecordNotFound = err{code: 20100, message: "数据异常"}

	// ErrNotLogin 未登录
	ErrNotLogin = err{code: 20201, message: "未登陆"}
	// ErrAuth  登录验证失败
	ErrAuth = err{code: 20202, message: "验证失败"}

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
	status  int    //http状态码 默认200
	code    int    // 错误码
	message string // 展示给用户看的
	err     error  // 保存内部错误信息
	callers string //保存调用文件名
}

// Status 修改返回的http状态码
func (e err) Status(status int) err {
	e.status = status
	return e
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
func status(status int) int {
	if status == 0 {
		return 200
	}
	return status
}

// DecodeDebugErr 解码错误, 获取 真是错误信息以及调用堆栈
func DecodeDebugErr(e error) (e2 error, caller string) {
	if e == nil {
		return nil, ""
	}

	if decodeErr, ok := e.(err); ok {
		return decodeErr.err, decodeErr.callers
	}
	return e, ""
}

var Debug = false

// DecodeErr 解码错误, 获取 httpStatus、 code 和 message
func DecodeErr(e error) (httpStatus int, code int, message string) {
	if e == nil {
		return status(OK.status), OK.Code(), OK.Message()
	}

	if decodeErr, ok := e.(err); ok {
		log.Println(decodeErr.code, decodeErr.message, decodeErr.err, decodeErr.callers)
		if Debug {
			return status(decodeErr.status), decodeErr.Code(), decodeErr.Error()
		}
		return status(decodeErr.status), decodeErr.Code(), decodeErr.Message()
	}
	if Debug {
		return status(ErrInternalServer.status), ErrInternalServer.Code(), e.Error()
	}
	return status(ErrInternalServer.status), ErrInternalServer.Code(), ErrInternalServer.Message()
}

func New(code int, message string) err {
	return err{
		code:    code,
		message: message,
	}
}
