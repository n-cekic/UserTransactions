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
	l.logger.Print("[INFO]" + fmt.Sprint(v...))
}

func (l *Log) Infof(s string, v ...any) {
	l.logger.Printf("INFO" + fmt.Sprintf(s, v...))
}

func (l *Log) Error(v ...any) {
	l.logger.Print("[ERROR]" + fmt.Sprint(v...))
}

func (l *Log) Errorf(s string, v ...any) {
	l.logger.Printf("ERROR" + fmt.Sprintf(s, v...))
}

func (l *Log) Fatal(v ...any) {
	l.logger.Fatal("[FATAL]" + fmt.Sprint(v...))
}

func (l *Log) Fatalf(s string, v ...any) {
	l.logger.Fatalf("FATAL" + fmt.Sprintf(s, v...))
}
