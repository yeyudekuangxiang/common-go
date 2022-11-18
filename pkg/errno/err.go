package errno

import (
	"errors"
	"fmt"
	"log"
	"mio/config"
	"mio/pkg/wxwork"
	"os"
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

	ErrCommon         = err{code: 10000, message: "系统繁忙,请稍后再试"}   // ErrCommon 通用错误
	ErrInternalServer = err{code: 10001, message: "内部服务器错误"}      // ErrInternalServer 系统错误
	ErrBind           = err{code: 10002, message: "参数错误"}         // ErrBind 参数错误
	ErrLimit          = err{code: 10003, message: "操作太频繁了、请稍后再试"} // ErrLimit 超出频率限制
	ErrTimeout        = err{code: 10004, message: "操作已超时"}        //超时

	ErrRecordNotFound  = err{code: 20100, message: "数据异常"}   // ErrRecordNotFound 数据库错误
	ErrChannelNotFound = err{code: 20101, message: "渠道数据异常"} // ErrChannelNotFound 渠道数据异常
	ErrExisting        = err{code: 20102, message: "数据已存在"}  // ErrExisting 数据已存在
	ErrCreate          = err{code: 20102, message: "保存失败"}   // ErrExisting 保存失败
	ErrUpdate          = err{code: 20102, message: "更新失败"}   // ErrExisting 更新失败
	ErrDelete          = err{code: 20102, message: "删除失败"}   // ErrExisting 删除失败

	ErrAuth       = err{code: 20201, message: "未登陆"}  // ErrAuth 未登录
	ErrValidation = err{code: 20202, message: "验证失败"} // ErrValidation 验证失败

	ErrUserNotFound       = err{code: 20301, message: "未查询到用户信息"} // ErrUserNotFound 未查询到用户信息
	ErrNotBindMobile      = err{code: 20303, message: "未授权手机号码"}  // ErrNotBindMobile 未绑定手机号
	ErrBindMobile         = err{code: 20304, message: "绑定手机号码失败"} // ErrBindMobile 绑定手机号时异常
	ErrBindRecordNotFound = err{code: 20305, message: "未找到绑定关系"}
	ErrNotFound           = err{code: 20305, message: "未找到资源"}

	ErrAdminNotFound = err{code: 20401, message: "未查询到管理员信息"} // ErrAdminNotFound 管理员错误 前缀204

	// 其他
	ErrCheckErr = err{code: 40001, message: "审核未通过"}
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
	_, f, l, _ := runtime.Caller(1)
	e.callers = f + ":" + strconv.Itoa(l)
	return e
}

// WithMessage 替换默认的提示
func (e err) WithMessage(message string) err {
	_, f, l, _ := runtime.Caller(1)
	e.callers = f + ":" + strconv.Itoa(l)
	e.message = message
	return e
}

// WithErrMessage 带上err message
func (e err) WithErrMessage(err string) err {
	e.err = errors.New(err)
	_, f, l, _ := runtime.Caller(1)
	e.callers = f + ":" + strconv.Itoa(l)
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

// DecodeErr 解码错误, 获取 httpStatus、 code 和 message
func DecodeErr(e error) (httpStatus int, code int, message string) {
	if e == nil {
		return status(OK.status), OK.Code(), OK.Message()
	}
	if decodeErr, ok := e.(err); ok {
		logerr(e, decodeErr.callers)
		if config.Config.App.Debug && decodeErr.err != nil {
			return status(decodeErr.status), decodeErr.Code(), decodeErr.err.Error()
		}
		return status(decodeErr.status), decodeErr.Code(), decodeErr.Message()
	}
	if _, ok := e.(fmt.Formatter); ok {
		logerr(e, fmt.Sprintf("%+v", e))
	} else {
		logerr(e, "")
	}

	if config.Config.App.Debug {
		return status(ErrInternalServer.status), ErrInternalServer.Code(), e.Error()
	}
	return status(ErrInternalServer.status), ErrInternalServer.Code(), ErrInternalServer.Error()

	//后面系统错误全面替换后使用下面的方式
	//return ErrInternalServer.Code(), ErrInternalServer.Message()
}
func logerr(err error, callers string) {
	if config.Config.App.Env != "prod" {
		return
	}
	go func() {
		sendErr := wxwork.SendRobotMessage(config.Constants.WxWorkBugRobotKey, wxwork.Markdown{
			Content: fmt.Sprintf("**容器:**%s \n\n**来源:**响应 \n\n**消息:**%+v \n\n**堆栈:**%+v", os.Getenv("HOSTNAME"), err.Error(), callers),
		})
		if sendErr != nil {
			log.Printf("推送异常到企业微信失败 %v %v", err, sendErr)
		}
	}()
}
