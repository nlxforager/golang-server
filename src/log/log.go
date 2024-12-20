package log

import (
	"fmt"
	"log/slog"
)

type logger struct {
	*slog.Logger
}

var Logger = slog.Default()

var UnknownLoggerType = fmt.Errorf("unknown logger type")

func Set(loggerType string) error {
	switch loggerType {
	case "TEST":
		Logger = newTestLogger()
	default:
		return UnknownLoggerType
	}
	return nil
}
