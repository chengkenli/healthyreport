package util

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

// logPath 日志文件路径
// logLevel 日志级别 debug/info/warn/error
// maxSize 单个文件大小,MB
// maxBackups 保存的文件个数
// maxAge 保存的天数， 没有的话不删除
// compress 压缩
// jsonFormat 是否输出为json格式
// shoowLine 显示代码行
// logInConsole 是否同时输出到控制台
var (
	logPath      string
	logLevel     string
	maxSize      int
	maxBackups   int
	maxAge       int
	compress     bool
	jsonFormat   bool
	showLine     bool
	logInConsole bool
)

func init() {
	logPath = fmt.Sprintf("%s/.%s.log", ".", filepath.Base(os.Args[0]))
	logLevel = "0"
	maxSize = 0
	maxBackups = 0
	maxAge = 0
	compress = false
	jsonFormat = false
	showLine = true
	logInConsole = true
}

func InitLogger() {
	hook := lumberjack.Logger{
		Filename:   logPath,    // 日志文件路径
		MaxSize:    maxSize,    // megabytes
		MaxBackups: maxBackups, // 最多保留300个备份
		Compress:   compress,   // 是否压缩 disabled by default
	}
	if maxAge > 0 {
		hook.MaxAge = maxAge // days
	}

	var syncer zapcore.WriteSyncer
	if logInConsole {
		syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	} else {
		syncer = zapcore.AddSync(&hook)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	var encoder zapcore.Encoder
	if jsonFormat {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 设置日志级别,debug可以打印出info,debug,warn；info级别可以打印warn，info；warn只能打印warn
	// debug->info->warn->error
	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	core := zapcore.NewCore(
		encoder,
		syncer,
		level,
	)

	Logger = zap.New(core)
	if showLine {
		Logger = Logger.WithOptions(zap.AddCaller())
	}
}

func Loggers() {
	//logger
	_, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}

	InitLogger()
	Logger.Info("log format is success.")
}
