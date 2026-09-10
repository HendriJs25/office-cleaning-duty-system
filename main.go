package main

import (
	"cleaning/cmd"
	"log/slog"
	"os"
	"time"
)

func main() {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		slog.Error("error loading timezone", "error", err)
		os.Exit(1)
	}

	time.Local = loc

	if err := cmd.Execute(); err != nil {
		slog.Error("execute failed", "error", err)
		os.Exit(1)
	}
}
