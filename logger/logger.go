package logger

import (
	"time"

	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	DPanic(args ...interface{})
	DPanicf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

type logger struct {
	log  *zap.Logger
	slog *zap.SugaredLogger
}

func New(level string) Logger {
	l := &logger{}
	l.initZap(level)
	return l
}

func (l *logger) initZap(level string) {
	cores := []zapcore.Core{}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(ParseLevel(level))

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeName = zapcore.FullNameEncoder
	encoderConfig.ConsoleSeparator = "\t"

	cores = append(cores, zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(colorable.NewColorableStdout()),
		atomicLevel,
	))

	zapOpts := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}
	if level == "debug" {
		zapOpts = append(zapOpts, zap.Development(), zap.AddStacktrace(zapcore.ErrorLevel))
	}

	l.log = zap.New(zapcore.NewTee(cores...), zapOpts...)
	l.slog = l.log.Sugar()

	defer l.Sync()
}

func (l *logger) Logger() *zap.Logger {
	return l.log
}

func (l *logger) SugaredLogger() *zap.SugaredLogger {
	return l.slog
}

func (l *logger) Sync() {
	if l.log != nil {
		l.log.Sync()
	}
	if l.slog != nil {
		l.slog.Sync()
	}
}

func (l *logger) Debug(args ...interface{}) {
	l.slog.Debug(args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.slog.Debugf(format, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.slog.Info(args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.slog.Infof(format, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.slog.Warn(args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.slog.Warnf(format, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.slog.Error(args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.slog.Errorf(format, args...)
}

func (l *logger) DPanic(args ...interface{}) {
	l.slog.DPanic(args...)
}

func (l *logger) DPanicf(format string, args ...interface{}) {
	l.slog.DPanicf(format, args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.slog.Panic(args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.slog.Panicf(format, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.slog.Fatal(args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.slog.Fatalf(format, args...)
}
