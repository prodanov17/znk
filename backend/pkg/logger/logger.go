package logger

import (
	"log"
	"os"

	"github.com/prodanov17/znk/internal/config"
)

type logger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	logLevel    int
}

const (
	INFO = iota
	WARN
	ERROR
)

var Log *logger = NewLogger()

func NewLogger() *logger {
	flags := log.LstdFlags
	var output *os.File
	var outputError *os.File
	var err error

	if config.Env.Env == "prod" {
		output, err = os.OpenFile("./storage/logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		outputError, err = os.OpenFile("./storage/logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		output = os.Stdout
		outputError = os.Stderr
	}

	return &logger{
		infoLogger:  log.New(output, "INFO: ", flags),
		warnLogger:  log.New(output, "WARN: ", flags),
		errorLogger: log.New(outputError, "ERROR: ", flags),
		logLevel:    INFO,
	}
}

func (l *logger) Info(v ...interface{}) {
	if l.logLevel <= INFO {
		l.infoLogger.Println(v...)
	}
}

func (l *logger) Warn(v ...interface{}) {
	if l.logLevel <= WARN {
		l.warnLogger.Println(v...)
	}
}

func (l *logger) Error(v ...interface{}) {
	l.errorLogger.Println(v...)
}

func (l *logger) Fatal(v ...interface{}) {
	l.errorLogger.Fatal(v...)
}

func (l *logger) SetLogLevel(level int) {
	l.logLevel = level
}

func (l *logger) Infof(format string, v ...interface{}) {
	if l.logLevel <= INFO {
		l.infoLogger.Printf(format, v...)
	}
}

func (l *logger) Warnf(format string, v ...interface{}) {
	if l.logLevel <= WARN {
		l.warnLogger.Printf(format, v...)
	}
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}
