package logger

import (
	"auth-template/internal/config"
	"log"
	"os"
)

type Logger struct {
	logger *log.Logger
	level  string
}

func NewLogger(cfg *config.Config) *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		level:  cfg.Log.Level,
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level == "debug" {
		l.logger.Printf("[DEBUG] "+format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.logger.Printf("[INFO] "+format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.logger.Printf("[WARN] "+format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.logger.Printf("[ERROR] "+format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.logger.Fatalf("[FATAL] "+format, v...)
}
