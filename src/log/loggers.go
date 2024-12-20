package log

import (
	"log/slog"
	"os"
)

const LevelSystem = slog.Level(-8)
const LevelSystemValue = "SYSTEM"

func newTestLogger() *slog.Logger {
	th := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: LevelSystem,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				if level == LevelSystem {
					a.Value = slog.StringValue(LevelSystemValue)
				}
			}
			return a
		},
	})
	return slog.New(th)
}
