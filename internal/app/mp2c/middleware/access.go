package middleware

import (
	"bufio"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/metric"
	"go.uber.org/zap"
	"io/ioutil"
	"strconv"
	"time"
)

const serverNamespace = "http_server"

var (
	metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "http server requests duration(ms).",
		Labels:    []string{"path"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "http server requests error count.",
		Labels:    []string{"path", "code"},
	})
)

func Access(logger *zap.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		bodyStr := "body too long"

		//记录长度小于5KB(5120个字节)的body
		if c.Request.ContentLength <= 1024*5 && c.Request.ContentLength >= 0 {
			body, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				logger.Error(err.Error())
			}
			c.Request.Body.Close()

			c.Request.Body = ioutil.NopCloser(bufio.NewReader(bytes.NewReader(body)))
			bodyStr = string(body)
		}

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			logger.Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("time", end.Format(timeFormat)),
				zap.Duration("latency", latency),
				zap.String("body", bodyStr),
			)
		}

		metricServerReqDur.Observe(int64(latency/time.Millisecond), path)
		metricServerReqCodeTotal.Inc(path, strconv.Itoa(c.Writer.Status()))
	}
}
