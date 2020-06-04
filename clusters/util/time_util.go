package util

import (
	"log"
	"time"
)

// TrackTime ...
func TrackTime(s string, startTime time.Time) {
	endTime := time.Now()
	timeSec := float64(endTime.Sub(startTime)) / float64(time.Second)

	log.Printf("%s took %.1f s", s, timeSec)
}
