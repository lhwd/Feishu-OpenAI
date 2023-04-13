package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

func LogInit() {
	logFile, err := os.OpenFile("./info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetPrefix("larkOpenAI-")
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
}

func LoggerMiddleware() gin.HandlerFunc {
	logClient := logrus.New()

	//logrus的输出
	src, err := os.OpenFile("./gin_info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("open logrus log err, err is ", err)
	}
	logClient.Out = src
	logClient.SetLevel(logrus.TraceLevel)
	apiLogPath := "api.log"
	logWriter, err := rotatelogs.New(
		apiLogPath+".%Y-%m-%d-%H-%M.log",
		rotatelogs.WithLinkName(apiLogPath),       // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.PanicLevel: logWriter,
		logrus.TraceLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.TextFormatter{})
	logClient.AddHook(lfHook)

	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		//执行时间
		latency := end.Sub(start)

		path := c.Request.URL.Path

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		logClient.Infof("| %3d | %13v | %15s | %s  %s |",
			statusCode,
			latency,
			clientIP,
			method, path,
		)
	}
}
