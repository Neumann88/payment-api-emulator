package loggin

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ILogger interface {
	Debugf(message string, args ...interface{})
	Debug(args ...interface{})
	Infof(message string, args ...interface{})
	Info(args ...interface{})
	Warnf(message string, args ...interface{})
	Warn(args ...interface{})
	Errorf(message string, args ...interface{})
	Error(args ...interface{})
	Fatalf(message string, args ...interface{})
	Fatal(args ...interface{})
}

type logger struct {
	logger ILogger
}

func NewLogger() *logger {
	return &logger{}
}

func (l *logger) Init(debug bool) {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "time"
	config.MessageKey = "message"
	config.EncodeTime = zapcore.TimeEncoderOfLayout("02-01-2006,15:04:05")
	config.EncodeCaller = zapcore.ShortCallerEncoder
	encoderJSON := zapcore.NewJSONEncoder(config)
	encoderConsole := zapcore.NewConsoleEncoder(config)

	err := os.MkdirAll("logs", 0755)
	if err != nil {
		log.Fatalln(err.Error())
	}

	logFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalln(err.Error())
	}

	var defaultLogLevel zapcore.Level
	if debug {
		defaultLogLevel = zapcore.DebugLevel
	} else {
		defaultLogLevel = zapcore.InfoLevel
	}

	core := zapcore.NewTee(
		zapcore.NewCore(encoderJSON, zapcore.AddSync(logFile), defaultLogLevel),
		zapcore.NewCore(encoderConsole, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)

	l.logger = zap.New(core, zap.AddCallerSkip(1), zap.AddCaller()).Sugar()
}

func (l *logger) Debugf(message string, args ...interface{}) {
	l.logger.Debugf(message, args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *logger) Infof(message string, args ...interface{}) {
	l.logger.Infof(message, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *logger) Warnf(message string, args ...interface{}) {
	l.logger.Warnf(message, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *logger) Errorf(message string, args ...interface{}) {
	l.logger.Errorf(message, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *logger) Fatalf(message string, args ...interface{}) {
	l.logger.Fatalf(message, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}
