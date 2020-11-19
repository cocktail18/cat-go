package cat

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warning(format string, args ...interface{})
	Error(format string, args ...interface{})
	SetOutput(w io.Writer)
}

type DefaultLogger struct {
	logger     *log.Logger
	mu         sync.Mutex
	currentDay int
}

func createLogger() *DefaultLogger {
	now := time.Now()

	var writer = getWriterByTime(now)

	return &DefaultLogger{
		logger:     log.New(writer, "", log.LstdFlags),
		mu:         sync.Mutex{},
		currentDay: now.Day(),
	}
}

func openLoggerFile(time time.Time) (*os.File, error) {
	year, month, day := time.Date()
	filename := fmt.Sprintf("%s/cat_%d_%02d_%02d.log", config.logDir, year, month, day)
	return os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
}

func getWriterByTime(time time.Time) io.Writer {
	if file, err := openLoggerFile(time); err != nil {
		log.Fatalf("Cannot open log file: %s, logs will be redirected to stdout", file.Name())
		return os.Stdout
	} else {
		log.Printf("Log has been redirected to the file: %s", file.Name())
		return file
	}
}

func (l *DefaultLogger) switchLogFile(time time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.currentDay == time.Day() {
		return
	}
	l.logger.SetOutput(getWriterByTime(time))
}

func (l *DefaultLogger) write(prefix, format string, args ...interface{}) {
	now := time.Now()

	if now.Day() != l.currentDay {
		l.switchLogFile(now)
	}
	l.logger.Printf(prefix+" "+format, args...)
}

func (l *DefaultLogger) Debug(format string, args ...interface{}) {
	l.write("[Debug]", format, args...)
}

func (l *DefaultLogger) Info(format string, args ...interface{}) {
	l.write("[Info]", format, args...)
}

func (l *DefaultLogger) Warning(format string, args ...interface{}) {
	l.write("[Warning]", format, args...)
}

func (l *DefaultLogger) Error(format string, args ...interface{}) {
	l.write("[Error]", format, args...)
}

func (l *DefaultLogger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

var logger Logger = createLogger()

