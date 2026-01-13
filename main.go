package main

import (
	"log/slog"
	"os"
)

func main() {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	log := slog.New(logHandler)

	log.Info("App started")

	
}
