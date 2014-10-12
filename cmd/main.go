package main

import (
	"github.com/rainerborene/cloudwatcher"
	"time"
)

func main() {
	// Parse configuration
	cloudwatcher.Config.Parse()

	// Collect metrics
	collector := cloudwatcher.NewCollector(5 * time.Minute)
	<-collector.Collect()
}
