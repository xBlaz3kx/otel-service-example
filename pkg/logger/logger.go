package logger

import (
	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"log/syslog"
)

type Logger struct {
	logger *log.Logger
	fields log.Fields
}

func NewLogger(address string) *Logger {
	logger := log.New()
	logger.SetLevel(log.DebugLevel)

	hook, err := lSyslog.NewSyslogHook(
		"tcp",
		address,
		syslog.LOG_DEBUG,
		"central-system",
	)
	if err == nil {
		logger.AddHook(hook)
	}

	return &Logger{
		logger: logger,
	}
}

func (l *Logger) SetFields(fields log.Fields) {
	l.fields = fields
}

func (l *Logger) WithFields(fields log.Fields) *log.Entry {
	return l.logger.WithFields(l.fields).WithFields(fields)
}

func (l *Logger) Get() *log.Entry {
	return l.logger.WithFields(l.fields)
}

func (l *Logger) WithTrace(traceId, spanId string) *log.Entry {
	return l.logger.WithFields(l.fields).WithFields(log.Fields{
		"traceId": traceId,
		"spanId":  spanId,
	})
}
