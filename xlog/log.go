package xlog

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	FormatText = "text"
	FormatJSON = "json"
)

var logger *zap.Logger

type LogConfig struct {
	LogPath    string // file path, if empty console will be used
	LogLevel   string // level
	Compress   bool   // compress or not
	MaxSize    int    // log size (MB)
	MaxAge     int    // lifecycle (day)
	MaxBackups int    // backup nums
	Format     string // text or json
}

func defaultOption() *LogConfig {
	return &LogConfig{
		LogPath:    "",
		LogLevel:   "debug",
		Compress:   true,
		MaxSize:    10,
		MaxAge:     90,
		MaxBackups: 20,
		Format:     FormatText,
	}
}

type Option func(cfg *LogConfig)

func getZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func newLogWriter(logCfg *LogConfig) io.Writer {
	if logCfg.LogPath == "" || logCfg.LogPath == "-" {
		return os.Stdout
	}
	return &lumberjack.Logger{
		Filename:   logCfg.LogPath,
		MaxSize:    logCfg.MaxSize, // default 100MB
		MaxAge:     logCfg.MaxAge,
		MaxBackups: logCfg.MaxBackups,
		LocalTime:  false,
		Compress:   logCfg.Compress,
	}
}

func newLoggerCore(logCfg *LogConfig) zapcore.Core {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		FunctionKey:    "func",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000Z0700"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	var encoder zapcore.Encoder
	if logCfg.Format == FormatJSON {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	return zapcore.NewCore(encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(newLogWriter(logCfg))),
		zap.NewAtomicLevelAt(getZapLevel(logCfg.LogLevel)))
}

func InitLog(opts ...Option) error {
	logCfg := defaultOption()
	for _, opt := range opts {
		opt(logCfg)
	}
	logger = zap.New(newLoggerCore(logCfg),
		//zap.AddCaller(),
		//zap.AddCallerSkip(1),
		//zap.AddStacktrace(getZapLevel(logCfg.LogLevel)),
		zap.Development())
	return nil
}

func GetLogger(opts ...zap.Option) *zap.Logger {
	return logger.WithOptions(opts...)
}

func Sync() error {
	return logger.Sync()
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
