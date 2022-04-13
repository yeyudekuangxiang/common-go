package errno

import (
	"fmt"
)

/*
状态码参考 https://segmentfault.com/a/1190000020525050
第一位表示错误级别, 1 为系统错误, 2 为普通错误
第二三位表示服务模块代码
第四五位表示具体错误代码
*/
var (
	OK = Err{code: 200, message: "OK"}

	// 系统错误, 前缀为 100
	InternalServerError = Err{code: 10001, message: "内部服务器错误"}
	ErrBind             = Err{code: 10002, message: "请求参数错误"}
	ErrTokenSign        = Err{code: 10003, message: "签名 jwt 时发生错误"}
	ErrEncrypt          = Err{code: 10004, message: "加密用户密码时发生错误"}

	// 数据库错误, 前缀为 201
	ErrDatabase = Err{code: 20100, message: "数据库错误"}
	ErrFill     = Err{code: 20101, message: "从数据库填充 struct 时发生错误"}

	// 认证错误, 前缀是 202
	ErrAuth         = Err{code: 20201, message: "未登陆"}
	ErrValidation   = Err{code: 20202, message: "验证失败"}
	ErrTokenInvalid = Err{code: 20203, message: "jwt 是无效的"}

	// 用户错误, 前缀为 203
	ErrUserNotFound      = Err{code: 20301, message: "用户没找到"}
	ErrPasswordIncorrect = Err{code: 20302, message: "密码错误"}
	ErrNotBindMobile     = Err{code: 20303, message: "未授权手机号码"}
)

func NewBindErr(err error) Err {
	return Err{
		code:    ErrBind.Code(),
		message: err.Error(),
		err:     err,
	}
}
func NewInternalServerError(err error) Err {
	return Err{
		code:    InternalServerError.Code(),
		message: err.Error(),
		err:     err,
	}
}
func NewAuthErr(err error) Err {
	return Err{
		code:    ErrAuth.Code(),
		message: err.Error(),
		err:     err,
	}
}

type IErr interface {
	Code() int
	Message() string
	error
}

// Err 定义错误
type Err struct {
	code    int    // 错误码
	message string // 展示给用户看的
	err     error  // 保存内部错误信息
}

func (err Err) Code() int {
	return err.code
}
func (err Err) Message() string {
	return err.message
}
func (err Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s", err.code, err.message)
}

// DecodeErr 解码错误, 获取 Code 和 Message
func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code(), OK.Message()
	}
	if e, ok := err.(IErr); ok {
		return e.Code(), e.Message()
	}
	return InternalServerError.Code(), err.Error()
}
