package errno

import (
	"fmt"
	"log"
)

/*
第一位表示错误级别, 1 为系统错误, 2 为普通错误
第二三位表示服务模块代码
第四五位表示具体错误代码
*/
var (
	OK = Errno{Code: 200, Message: "OK"}

	// 系统错误, 前缀为 100
	InternalServerError = Errno{Code: 10001, Message: "内部服务器错误"}
	ErrBind             = Errno{Code: 10002, Message: "请求参数错误"}
	ErrTokenSign        = Errno{Code: 10003, Message: "签名 jwt 时发生错误"}
	ErrEncrypt          = Errno{Code: 10004, Message: "加密用户密码时发生错误"}

	// 数据库错误, 前缀为 201
	ErrDatabase = Errno{Code: 20100, Message: "数据库错误"}
	ErrFill     = Errno{Code: 20101, Message: "从数据库填充 struct 时发生错误"}

	// 认证错误, 前缀是 202
	ErrAuth         = Errno{Code: 20201, Message: "未登陆"}
	ErrValidation   = Errno{Code: 20202, Message: "验证失败"}
	ErrTokenInvalid = Errno{Code: 20203, Message: "jwt 是无效的"}

	// 用户错误, 前缀为 203
	ErrUserNotFound      = Errno{Code: 20301, Message: "用户没找到"}
	ErrPasswordIncorrect = Errno{Code: 20302, Message: "密码错误"}
)

// 定义错误码
type Errno struct {
	Code    int
	Message string
}

func (err Errno) Error() string {
	return err.Message
}

// 定义错误
type Err struct {
	Code    int    // 错误码
	Message string // 展示给用户看的
	Errord  error  // 保存内部错误信息
}

func (err Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.Errord)
}

// 使用 错误码 和 error 创建新的 错误
func New(errno Errno, err error) Err {
	return Err{
		Code:    errno.Code,
		Message: err.Error(),
		Errord:  err,
	}
}

//使用 error 创建新的 绑定错误
func NewBindErr(err error) Err {
	return Err{
		Code:    ErrBind.Code,
		Message: err.Error(),
		Errord:  err,
	}
}

//使用 错误码 和 新的错误信息 创建新的 错误
func NewWithMessage(errno Errno, message string) Err {
	return Err{
		Code:    errno.Code,
		Message: message,
		Errord:  nil,
	}
}

func DefaultErr(err error) Err {
	return Err{
		Code:    InternalServerError.Code,
		Message: err.Error(),
		Errord:  err,
	}
}

// 解码错误, 获取 Code 和 Message
func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}
	switch typed := err.(type) {
	case Err:
		log.Printf("%+v", err)
		return typed.Code, typed.Message
	case Errno:
		return typed.Code, typed.Message
	default:
	}

	return InternalServerError.Code, err.Error()
}
