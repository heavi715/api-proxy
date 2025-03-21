package logger

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"sync"
	"time"
)

var logger *logrus.Logger
var once sync.Once

func init() {
	//logger = NewLogger()
	once.Do(func() {
		logger = NewLogger()
	})
}

func Info(args ...interface{}) {
	logger.Infoln(args...)
}

// Error
func Error(args ...interface{}) {
	logger.Errorln(args...)
}

// Warning 警告
func Warning(args ...interface{}) {
	logger.Warningln(args...)
}

// DeBug debug
func DeBug(args ...interface{}) {
	logger.Debugln(args...)
}

func NewLogger() *logrus.Logger {
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println(err.Error())
	}
	logFileName := now.Format("2006-01-02") + ".log"
	//日志文件
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			fmt.Println(err.Error())
		}
	}
	//写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	//实例化
	logger := logrus.New()

	//设置输出
	logger.Out = src

	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return logger
}

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RequestLog() gin.HandlerFunc {

	return func(c *gin.Context) {

		crw := &CustomResponseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = crw
		reqBody, _ := c.GetRawData()
		// 请求包体写回。
		if len(reqBody) > 0 {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 请求IP
		clientIP := c.ClientIP()

		//日志格式
		logger.Infof("| %15s | %s | %s | %s |",
			clientIP,
			reqMethod,
			reqUri,
			reqBody,
		)

		// 处理请求
		c.Next()

	}
}

func ResponseLog() gin.HandlerFunc {
	//logger := Logger()
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		crw := &CustomResponseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = crw
		reqBody, _ := c.GetRawData()
		// 请求包体写回。
		if len(reqBody) > 0 {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		respBody := string(crw.body.Bytes())

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		//日志格式
		logger.Infof("| %3d | %13v | %15s | %s | %s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
			reqBody,
			respBody,
		)
	}
}
