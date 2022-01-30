package log

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/azi-v/ladon-api/config"
	"github.com/ory/ladon"
)

// log config
func Init(cfg *config.LogConfig) {
	// logger初始化
	DefaultLogger = NewLogger(cfg)
}

// A Level is a logging priority. Higher levels are more important.
type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel

	_minLevel = DebugLevel
	_maxLevel = FatalLevel
)

// String returns a lower-case ASCII representation of the log level.
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case PanicLevel:
		return "panic"
	case FatalLevel:
		return "fatal"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

// CapitalString returns an all-caps ASCII representation of the log level.
func (l Level) CapitalString() string {
	// Printing levels in all-caps is common enough that we should export this
	// functionality.
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case PanicLevel:
		return "PANIC"
	case FatalLevel:
		return "FATAL"
	default:
		return fmt.Sprintf("LEVEL(%d)", l)
	}
}

type ComonLogStruct struct {
	LogLevel    Level      `json:"log_level,omitempty"`
	Datetime    *time.Time `json:"datetime,omitempty"`
	ServiceName string     `json:"service_name,omitempty"`
	CallerKey   string     `json:"caller_key,omitempty"`
	FunctionKey string     `json:"function_key,omitempty"`
	Msg         string     `json:"msg,omitempty"`
	Context     struct {
		Input  map[string]interface{} `json:"input,omitempty"`
		Output map[string]interface{} `json:"output,omitempty"`
	} `json:"context,omitempty"`
	KeyData map[string]interface{} `json:"key_data,omitempty"`
}

func (cl *ComonLogStruct)Marshal() string {
	// 字符串拼接
	return ""
}

var DefaultLogger *Logger

// 所有日志都输出到标准输出与文件；
type Logger struct {
	*log.Logger
}

func NewLogger(cfg *config.LogConfig) *Logger {
	var out io.Writer
	var prefix string
	var flag int
	return &Logger{
		Logger: log.New(out, prefix, flag),
	}
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
	DefaultLogger.Info(ctx, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	DefaultLogger.Debug(ctx, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	DefaultLogger.Error(ctx, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	DefaultLogger.Warn(ctx, args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	DefaultLogger.Panic(ctx, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	DefaultLogger.Fatal(ctx, args...)
}
