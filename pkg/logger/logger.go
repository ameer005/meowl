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
		Level: slog.LevelDebug, // or Info
	})
	Logger = slog.New(handler)
}
