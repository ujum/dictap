package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/ujum/dictap/internal/config"
	"os"
)

type Logger interface {
	Debug(...interface{})
	Debugf(format string, args ...interface{})
	Info(...interface{})
	Infof(format string, args ...interface{})
	Warn(...interface{})
	Warnf(format string, args ...interface{})
	Error(...interface{})
	Errorf(format string, args ...interface{})
	Fatal(...interface{})
	Fatalf(format string, args ...interface{})
	Print(...interface{})
	Println(...interface{})
}

type LogrusAdaper struct {
	logger *logrus.Logger
}

func NewLogrus(cfg *config.Config) *LogrusAdaper {
	level, err := logrus.ParseLevel(cfg.Logger.Level)
	if err != nil {
		logrus.Warnf("can't parse '%s' log level, using 'info'", cfg.Logger.Level)
		level = logrus.InfoLevel
	}
	return &LogrusAdaper{
		logger: &logrus.Logger{
			Out:          os.Stdout,
			Formatter:    &logrus.JSONFormatter{},
			Hooks:        make(logrus.LevelHooks),
			Level:        level,
			ExitFunc:     os.Exit,
			ReportCaller: false,
		},
	}
}

func (adapter *LogrusAdaper) Debug(args ...interface{}) {
	adapter.logger.Debug(args...)
}

func (adapter *LogrusAdaper) Debugf(format string, args ...interface{}) {
	adapter.logger.Debugf(format, args...)
}

func (adapter *LogrusAdaper) Info(args ...interface{}) {
	adapter.logger.Info(args...)
}

func (adapter *LogrusAdaper) Infof(format string, args ...interface{}) {
	adapter.logger.Infof(format, args...)
}

func (adapter *LogrusAdaper) Warn(args ...interface{}) {
	adapter.logger.Warn(args...)
}

func (adapter *LogrusAdaper) Warnf(format string, args ...interface{}) {
	adapter.logger.Warnf(format, args...)
}

func (adapter *LogrusAdaper) Error(args ...interface{}) {
	adapter.logger.Error(args...)
}

func (adapter *LogrusAdaper) Errorf(format string, args ...interface{}) {
	adapter.logger.Errorf(format, args...)
}

func (adapter *LogrusAdaper) Fatal(args ...interface{}) {
	adapter.logger.Fatal(args...)
}

func (adapter *LogrusAdaper) Fatalf(format string, args ...interface{}) {
	adapter.logger.Fatalf(format, args...)
}

func (adapter *LogrusAdaper) Print(args ...interface{}) {
	adapter.logger.Print(args...)
}

func (adapter *LogrusAdaper) Println(args ...interface{}) {
	adapter.logger.Println(args...)
}
