package log

import (
	"github.com/astaxie/beego/logs"
	"log"
)

var (
	l          *log.Logger
	switch_log = true
)

func init() {
	l = logs.GetLogger()
	logs.SetLogger("console")
	logs.EnableFuncCallDepth(true)
}

func Println(p ...interface{}) {
	if switch_log {
		l.Println(p...)
	}
}

func Debug(p ...interface{}) {
	if switch_log {
		logs.Debug(p)
	}
}

func Info(p ...interface{}) {
	if switch_log {
		logs.Info(p)
	}
}

func Warn(p ...interface{}) {
	if switch_log {
		logs.Warn(p)
	}
}

func Error(p ...interface{}) {
	if switch_log {
		logs.Error(p)
	}
}

func Critical(p ...interface{}) {
	if switch_log {
		logs.Critical(p)
	}
}
