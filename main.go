package main

import (
	"cleaning/cmd"
	"log"
	"time"
)

func main() {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("failed to load timezone: %v", err)
	}

	time.Local = loc

	if err := cmd.Execute(); err != nil {
		log.Fatalf("execute failed: %v", err)
	}
}
