package logger

import (
	"log"
	"os"
)

// Logger interface for logging
type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatalf(format string, v ...interface{})
}

type logger struct {
	*log.Logger
}

// New creates a new logger
func New(level string) Logger {
	return &logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *logger) Printf(format string, v ...interface{}) {
	l.Logger.Printf(format, v...)
}

func (l *logger) Println(v ...interface{}) {
	l.Logger.Println(v...)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatalf(format, v...)
}
