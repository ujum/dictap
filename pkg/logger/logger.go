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

type logrusAdaper struct {
	logger logrus.Logger
}

func New(cfg *config.Config) Logger {
	level, err := logrus.ParseLevel(cfg.Logger.Level)
	if err != nil {
		logrus.Warnf("can't parse '%s' log level, using 'info'", cfg.Logger.Level)
		level = logrus.InfoLevel
	}
	return &logrusAdaper{
		logger: logrus.Logger{
			Out:          os.Stdout,
			Formatter:    &logrus.JSONFormatter{},
			Hooks:        make(logrus.LevelHooks),
			Level:        level,
			ExitFunc:     os.Exit,
			ReportCaller: false,
		},
	}
}

func (adapter *logrusAdaper) Debug(args ...interface{}) {
	adapter.logger.Debug(args)
}

func (adapter *logrusAdaper) Debugf(format string, args ...interface{}) {
	adapter.logger.Debugf(format, args)
}

func (adapter *logrusAdaper) Info(args ...interface{}) {
	adapter.logger.Info(args)
}

func (adapter *logrusAdaper) Infof(format string, args ...interface{}) {
	adapter.logger.Infof(format, args)
}

func (adapter *logrusAdaper) Warn(args ...interface{}) {
	adapter.logger.Warn(args)
}

func (adapter *logrusAdaper) Warnf(format string, args ...interface{}) {
	adapter.logger.Warnf(format, args)
}

func (adapter *logrusAdaper) Error(args ...interface{}) {
	adapter.logger.Error(args)
}

func (adapter *logrusAdaper) Errorf(format string, args ...interface{}) {
	adapter.logger.Errorf(format, args)
}

func (adapter *logrusAdaper) Fatal(args ...interface{}) {
	adapter.logger.Fatal(args)
}

func (adapter *logrusAdaper) Fatalf(format string, args ...interface{}) {
	adapter.logger.Fatalf(format, args)
}

func (adapter *logrusAdaper) Print(args ...interface{}) {
	adapter.logger.Print(args)
}

func (adapter *logrusAdaper) Println(args ...interface{}) {
	adapter.logger.Println(args)
}
