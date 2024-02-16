package logging

import "log"

// var Logger *log.Logger

type Log struct {
	logger *log.Logger
}

var Logger Log

func init() {
	Logger.logger = log.Default()
}

func (l *Log) Printf(s string, v ...any) {
	l.logger.Printf(s, v...)
}

func (l *Log) Print(v ...any) {
	l.logger.Print(v...)
}

func (l *Log) Println(v ...any) {
	l.logger.Println(v...)
}

func (l *Log) Fatalf(s string, v ...any) {
	l.logger.Printf(s, v...)
}
