// Copyright (c) 2023, Viet Tran, 200Lab Team.

package sctx

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"runtime"
	"strings"
)

type Fields logrus.Fields

type Logger interface {
	Print(args ...interface{})
	Debug(...interface{})
	Debugln(...interface{})
	Debugf(string, ...interface{})

	Info(...interface{})
	Infoln(...interface{})
	Infof(string, ...interface{})

	Warn(...interface{})
	Warnln(...interface{})
	Warnf(string, ...interface{})

	Error(...interface{})
	Errorln(...interface{})
	Errorf(string, ...interface{})

	Fatal(...interface{})
	Fatalln(...interface{})
	Fatalf(string, ...interface{})

	Panic(...interface{})
	Panicln(...interface{})
	Panicf(string, ...interface{})

	With(key string, value interface{}) Logger
	Withs(Fields) Logger
	// add source field to log
	WithSrc() Logger
	GetLevel() string
}

type logger struct {
	*logrus.Entry
}

func (l *logger) GetLevel() string {
	return l.Entry.Logger.Level.String()
}

func (l *logger) debugSrc() *logrus.Entry {

	if _, ok := l.Entry.Data["source"]; ok {
		return l.Entry
	}

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	return l.Entry.WithField("source", fmt.Sprintf("%s:%d", file, line))
}

func (l *logger) Debug(args ...interface{}) {
	if l.Entry.Logger.Level >= logrus.DebugLevel {
		l.debugSrc().Debug(args...)
	}
}

func (l *logger) Debugln(args ...interface{}) {
	if l.Entry.Logger.Level >= logrus.DebugLevel {
		l.debugSrc().Debugln(args...)
	}
}

func (l *logger) Debugf(format string, args ...interface{}) {
	if l.Entry.Logger.Level >= logrus.DebugLevel {
		l.debugSrc().Debugf(format, args...)
	}
}

func (l *logger) Print(args ...interface{}) {
	if l.Entry.Logger.Level >= logrus.DebugLevel {
		l.debugSrc().Debug(args)
	}
}

func (l *logger) With(key string, value interface{}) Logger {
	return &logger{l.Entry.WithField(key, value)}
}

func (l *logger) Withs(fields Fields) Logger {
	return &logger{l.Entry.WithFields(logrus.Fields(fields))}
}

func (l *logger) WithSrc() Logger {
	return &logger{l.debugSrc()}
}

func mustParseLevel(level string) logrus.Level {
	lv, err := logrus.ParseLevel(level)

	if err != nil {
		log.Fatal(err.Error())
	}

	return lv
}

var (
	defaultLogger = newAppLogger(&Config{
		BasePrefix:   "core",
		DefaultLevel: "trace",
	})
)

type AppLogger interface {
	GetLogger(prefix string) Logger
}

func GlobalLogger() AppLogger {
	return defaultLogger
}

type Config struct {
	DefaultLevel string
	BasePrefix   string
}

type appLogger struct {
	logger   *logrus.Logger
	cfg      Config
	logLevel string
}

func newAppLogger(config *Config) *appLogger {
	if config == nil {
		config = &Config{}
	}

	if config.DefaultLevel == "" {
		config.DefaultLevel = "info"
	}

	logger := logrus.New()
	logger.Formatter = logrus.Formatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return &appLogger{
		logger:   logger,
		cfg:      *config,
		logLevel: config.DefaultLevel,
	}
}

func (al *appLogger) GetLogger(prefix string) Logger {
	var entry *logrus.Entry

	prefix = al.cfg.BasePrefix + "." + prefix
	prefix = strings.Trim(prefix, ".")
	if prefix == "" {
		entry = logrus.NewEntry(al.logger)
	} else {
		entry = al.logger.WithField("prefix", prefix)
	}

	return &logger{entry}
}

func (*appLogger) ID() string {
	return "logger"
}

func (al *appLogger) InitFlags() {
	flag.StringVar(&al.logLevel, "log-level", al.cfg.DefaultLevel, "Log level: panic | fatal | error | warn | info | debug | trace")
}

func (al *appLogger) Activate(_ ServiceContext) error {
	lv := mustParseLevel(al.logLevel)
	al.logger.SetLevel(lv)

	return nil
}

func (al *appLogger) Stop() error {
	return nil
}
