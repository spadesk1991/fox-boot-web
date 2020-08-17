package webErrors

import (
	"net/http"
)

// 错误处理的结构体
type Error struct {
	StatusCode int         `json:"-"`
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Result     interface{} `json:"result"`
}

func (e *Error) Error() string {
	return e.Msg
}

func NewError(statusCode, Code int, msg string) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       Code,
		Msg:        msg,
	}
}

var (
	Success     = NewError(http.StatusOK, 0, "success")
	ServerError = NewError(http.StatusInternalServerError, 200500, "系统异常，请稍后重试!")
	NotFound    = NewError(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
)

func BadRequestError(message string) *Error {
	return NewError(http.StatusBadRequest, 200400, message)
}

func UnauthorizedError() *Error {
	return NewError(http.StatusUnauthorized, 200401, "请登陆")
}
