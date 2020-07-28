package middleware

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"strings"
	"sync/atomic"
	"time"

	"github.com/spadesk1991/fox-boot-web/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	requestID uint64
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Log only when path is not being skipped
		var skip map[string]struct{}

		if _, ok := skip[path]; !ok {
			param := MyLogFormatterParams{
				LogFormatterParams: gin.LogFormatterParams{
					Request: c.Request,
					Keys:    c.Keys,
				},
			}

			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

			param.BodySize = c.Writer.Size()
			// body 参数
			reader := c.Request.Body

			if c.Request.Header.Get("Content-Encoding") == "gzip" {
				var err error
				reader, err = gzip.NewReader(c.Request.Body)
				if err != nil {
					logrus.Errorln(err)
				}
			}
			bt, e := ioutil.ReadAll(reader)
			if e != nil {
				logrus.Error(e)
			}
			c.Request.Body.Close()
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bt)) // 重新设置body
			param.Body = strings.ReplaceAll(string(bt), "\n", "")

			if raw != "" {
				path = path + "?" + raw
			}

			param.Path = path

			atomic.AddUint64(&requestID, 1)
			logger.LoggerClient().Infof("[WEB|  %v| %d | %s | %3d | %13v | %15s | %-7s  %#v body=%s\n%s",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				requestID,
				c.Request.Header.Get("User-Id"),
				param.StatusCode,
				param.Latency,
				param.ClientIP,
				param.Method,
				param.Path,
				param.Body,
				param.ErrorMessage,
			)
		}
		// Process request
		c.Next()
	}
}

type MyLogFormatterParams struct {
	gin.LogFormatterParams
	Body string
}
