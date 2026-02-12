package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	std              *log.Logger
	logID            string
	sessionID        string
	enableRequestIDs bool
}

func New() *Logger {
	return &Logger{
		std:              log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds),
		enableRequestIDs: false,
	}
}

func NewWithConfig(enableRequestIDs bool) *Logger {
	return &Logger{
		std:              log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds),
		enableRequestIDs: enableRequestIDs,
	}
}

func (l *Logger) WithContext(logID, sessionID string) *Logger {
	return &Logger{
		std:              l.std,
		logID:            logID,
		sessionID:        sessionID,
		enableRequestIDs: l.enableRequestIDs,
	}
}

func (l *Logger) prefix() string {
	if !l.enableRequestIDs {
		return ""
	}
	if l.logID == "" && l.sessionID == "" {
		return ""
	}

	return fmt.Sprintf("[log_id=%s session_id=%s] ", l.logID, l.sessionID)
}

func (l *Logger) Info(msg string) {
	l.std.Println("INFO:", l.prefix()+msg)
}

func (l *Logger) Error(msg string) {
	l.std.Println("ERROR:", l.prefix()+msg)
}
