package main

import (
	"time"
)

func main() {
	// Parse configuration
	Config.Parse()

	// Collect metrics
	collector := NewCollector(5 * time.Minute)
	<-collector.Run()
}
