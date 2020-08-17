package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/spadesk1991/fox-boot-web/webErrors"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
				var Err *webErrors.Error
				if e, ok := err.(*webErrors.Error); ok {
					Err = e
				} else if e, ok := err.(error); ok {
					logrus.Errorf("%v\n", err)
					Err = webErrors.BadRequestError(e.Error())
				} else {
					logrus.Errorf("%v\n", err)
					Err = webErrors.ServerError
				}
				// 记录一个错误的日志
				c.JSON(Err.StatusCode, Err)
				c.Abort()
			}
		}()
		c.Next()
	}
}

// 404处理
func HandleNotFound(c *gin.Context) {
	err := webErrors.NotFound
	c.JSON(err.StatusCode, err)
	return
}
