package initialize

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"server/global"
)

func InitLogger() {
	cfg := global.Config.Log

	// 解析日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 编码器配置（ELK 友好的字段命名）
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 文件输出始终使用 JSON 格式，方便 Filebeat/ELK 采集
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// 控制台输出按配置选择格式
	var consoleEncoder zapcore.Encoder
	if cfg.Format == "json" {
		consoleEncoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		consoleEncoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 应用默认值
	if cfg.Directory == "" {
		cfg.Directory = "./logs"
	}
	if cfg.Filename == "" {
		cfg.Filename = "app.log"
	}
	if cfg.MaxSize <= 0 {
		cfg.MaxSize = 100
	}
	if cfg.MaxBackups <= 0 {
		cfg.MaxBackups = 5
	}
	if cfg.MaxAge <= 0 {
		cfg.MaxAge = 30
	}

	// 构建输出 cores
	var cores []zapcore.Core

	// 文件输出（带 lumberjack 轮转）
	logFile := filepath.Join(cfg.Directory, cfg.Filename)
	fileWriter := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	cores = append(cores, zapcore.NewCore(jsonEncoder, zapcore.AddSync(fileWriter), level))

	// stdout 输出（默认开启）
	if cfg.Stdout {
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level))
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	global.Log = logger.Sugar()
}
