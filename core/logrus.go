package core

import (
	"Smart_delivery_locker/global"
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

// 日志颜色
const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

type LogFormatter struct{}

// Format 格式化日志 实现Formatter
func (t LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	//根据不同level展示颜色
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	log := global.Config.Logger

	//自定义日期格式
	timesTamp := entry.Time.Format("2006-01-02 15:04:05")
	if entry.HasCaller() {
		//自定义文件路径
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		//自定义输出路径
		_, _ = fmt.Fprintf(b, "%s[%s] [%dm][%s]\x1b[0m %s %s %s\n", log.Prefix, timesTamp, levelColor, entry.Level, funcVal, fileVal, entry.Message)
	} else {
		_, _ = fmt.Fprintf(b, "%s[%s] [%dm][%s]\x1b[0m %s\n", log.Prefix, timesTamp, levelColor, entry.Level, entry.Message)
	}
	return b.Bytes(), nil
}

func InitLogger() *logrus.Logger {
	mLog := logrus.New()                                //创建一个实例
	mLog.SetOutput(os.Stdout)                           //设置输出类型
	mLog.SetReportCaller(global.Config.Logger.ShowLine) //开启返回函数名 和 行号
	mLog.SetFormatter(&LogFormatter{})                  //设置自己定义的LogFormatter
	level, err := logrus.ParseLevel(global.Config.Logger.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	mLog.SetLevel(level) //设置最低的Level
	InitDefaultLogger()  //调用全局log
	return mLog
}

func InitDefaultLogger() {
	//全局log
	logrus.SetOutput(os.Stdout) //设置输出类型
	logrus.SetReportCaller(global.Config.Logger.ShowLine)
	logrus.SetFormatter(&LogFormatter{})
	level, err := logrus.ParseLevel(global.Config.Logger.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level) //设置最低的Level
}
