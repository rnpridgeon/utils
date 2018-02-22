package utils

import (
	"time"
	"log"
)

// TODO: create a simple stats library for ingesting/storing time tracking data
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("INFO: %s took %s", name, elapsed)
}