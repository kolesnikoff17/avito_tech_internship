package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

// Interface -.
type Interface interface {
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warnf(format string, args ...interface{})
	Warn(args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
}

// Logger is a wrap around different loggers, implements Interface
type Logger struct {
	l *zap.SugaredLogger
}

// New is a Logger constructor
func New(level string) (*Logger, error) {
	var l zapcore.Level

	switch strings.ToLower(level) {
	case "error":
		l = zapcore.ErrorLevel
	case "warn":
		l = zapcore.WarnLevel
	case "info":
		l = zapcore.InfoLevel
	case "debug":
		l = zapcore.DebugLevel
	default:
		l = zapcore.InfoLevel
	}

	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("Jan 02 15:04:05.000000000")
	config.Level = zap.NewAtomicLevelAt(l)

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	sugar := logger.Sugar()

	return &Logger{
		l: sugar,
	}, nil
}

// Errorf -.
func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.l.Errorf(format, args)
}

// Error -.
func (logger *Logger) Error(args ...interface{}) {
	logger.l.Error(args)
}

// Fatalf -.
func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.l.Fatalf(format, args)
}

// Fatal -.
func (logger *Logger) Fatal(args ...interface{}) {
	logger.l.Fatal(args)
}

// Infof -.
func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.l.Infof(format, args)
}

// Info -.
func (logger *Logger) Info(args ...interface{}) {
	logger.l.Info(args)
}

// Warnf -.
func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.l.Warnf(format, args)
}

// Warn -.
func (logger *Logger) Warn(args ...interface{}) {
	logger.l.Warn(args)
}

// Debugf -.
func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.l.Debugf(format, args)
}

// Debug -.
func (logger *Logger) Debug(args ...interface{}) {
	logger.l.Debug(args)
}
