package logging

import (
	"fmt"
	"log"
)

type Log struct {
	logger *log.Logger
}

var Logger Log

func init() {
	Logger.logger = log.Default()
}

func (l *Log) Info(v ...any) {
	l.logger.Print("[INFO]\t" + fmt.Sprint(v...))
}

func (l *Log) Infof(s string, v ...any) {
	l.logger.Printf("[INFO]\t" + fmt.Sprintf(s, v...))
}

func (l *Log) Error(v ...any) {
	l.logger.Print("[ERROR]\t" + fmt.Sprint(v...))
}

func (l *Log) Errorf(s string, v ...any) {
	l.logger.Printf("[ERROR]\t" + fmt.Sprintf(s, v...))
}

func (l *Log) Fatal(v ...any) {
	l.logger.Fatal("[FATAL]\t" + fmt.Sprint(v...))
}

func (l *Log) Fatalf(s string, v ...any) {
	l.logger.Fatalf("[FATAL]\t" + fmt.Sprintf(s, v...))
}
