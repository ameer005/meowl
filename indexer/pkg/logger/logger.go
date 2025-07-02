package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

var (
	Logger *slog.Logger
)

func Init() {
	handler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "15:04:05", // HH:MM:SS
		AddSource:  true,       // adds file:line of the log call
	})
	Logger = slog.New(handler)
}
