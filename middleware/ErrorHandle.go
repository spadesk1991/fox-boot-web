package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// 错误处理的结构体
type Error struct {
	StatusCode int         `json:"-"`
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Result     interface{} `json:"result"`
}

var (
	Success     = NewError(http.StatusOK, 0, "success")
	ServerError = NewError(http.StatusInternalServerError, 200500, "系统异常，请稍后重试!")
	NotFound    = NewError(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
)

func OtherError(message string) *Error {
	return NewError(http.StatusBadRequest, 100400, message)
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

// 404处理
func HandleNotFound(c *gin.Context) {
	err := NotFound
	c.JSON(err.StatusCode, err)
	return
}

func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
				var Err *Error
				if e, ok := err.(*Error); ok {
					Err = e
				} else if e, ok := err.(error); ok {
					logrus.Errorf("%v\n", err)
					Err = OtherError(e.Error())
				} else {
					logrus.Errorf("%v\n", err)
					Err = ServerError
				}
				// 记录一个错误的日志
				c.JSON(Err.StatusCode, Err)
				c.Abort()
			}
		}()
		c.Next()
	}
}
