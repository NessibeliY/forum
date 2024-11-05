package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"01.alem.school/git/nyeltay/forum/conf"
)

type Logger struct {
	*log.Logger
}

func Setup(config *conf.Config) (*Logger, error) {
	file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return nil, fmt.Errorf("open logfile: %w", err)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)

	l := &Logger{
		Logger: log.New(multiWriter, "", log.Ldate|log.Ltime|log.Lshortfile),
	}

	return l, nil
}

func (l *Logger) logWithCallerInfo(prefix string, v ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		l.SetPrefix(fmt.Sprintf("%s %s:%d ", prefix, file, line))
	} else {
		l.SetPrefix(prefix)
	}
	l.Logger.Println(v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logWithCallerInfo("FATAL:", fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.logWithCallerInfo("FATAL:", v...)
	os.Exit(1)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.logWithCallerInfo("INFO:", fmt.Sprintf(format, v...))
}

func (l *Logger) Info(v ...interface{}) {
	l.logWithCallerInfo("INFO:", v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logWithCallerInfo("ERROR:", fmt.Sprintf(format, v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.logWithCallerInfo("ERROR:", v...)
}
