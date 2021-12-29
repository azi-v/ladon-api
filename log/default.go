package log

import (
	"context"

	"git.ymt360.com/usercenter/ymt-ladon/config"
	"github.com/ory/ladon"
)

// log config
func Init() {
	// apollo配置拉取
	// logger初始化
	//
}

var DefaultLogger *Logger

type Logger struct{}

func NewLogger(cfg *config.LogConfig) *Logger {
	return &Logger{}
}

func (l *Logger) LogRejectedAccessRequest(request *ladon.Request, pool ladon.Policies, deciders ladon.Policies) {
}

func (l *Logger) LogGrantedAccessRequest(request *ladon.Request, pool ladon.Policies, deciders ladon.Policies) {
}

func (l *Logger) Info(ctx context.Context, args ...interface{}) {

}
func (l *Logger) Debug(ctx context.Context, args ...interface{}) {

}
func (l *Logger) Error(ctx context.Context, args ...interface{}) {

}
func (l *Logger) Warn(ctx context.Context, args ...interface{}) {

}
func (l *Logger) Panic(ctx context.Context, args ...interface{}) {

}
func (l *Logger) Fatal(ctx context.Context, args ...interface{}) {

}

func Info(ctx context.Context, args ...interface{}) {

}
func Debug(ctx context.Context, args ...interface{}) {

}
func Error(ctx context.Context, args ...interface{}) {

}
func Warn(ctx context.Context, args ...interface{}) {

}
func Panic(ctx context.Context, args ...interface{}) {

}
func Fatal(ctx context.Context, args ...interface{}) {

}
