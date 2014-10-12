package main

import (
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/cloudwatch"
	"log"
	"time"
)

type collector struct {
	server   *cloudwatch.CloudWatch
	duration time.Duration
}

func NewCollector(duration time.Duration) *collector {
	// AWS Authentication
	auth, err := aws.GetAuth("", "", "", time.Now())
	check(err)

	// Initialize CloudWatch
	region := aws.Regions[aws.InstanceRegion()]
	server, err := cloudwatch.NewCloudWatch(auth, region.CloudWatchServicepoint)
	check(err)

	return &collector{server: server, duration: duration}
}

func (c *collector) PutMetric(datum []cloudwatch.MetricDatum) {
	_, err := c.server.PutMetricDataNamespace(datum, Config.Namespace)

	check(err)

	for m := range datum {
		log.Printf("%s: %f (%s)\n", datum[m].MetricName, datum[m].Value, datum[m].Unit)
	}
}

func (c *collector) Run() chan bool {
	ticker := time.NewTicker(c.duration)
	stop := make(chan bool, 1)

	go func() {
		for {
			select {
			case <-ticker.C:
				c.PutMetric(GetMemoryDatum())
				c.PutMetric(GetFileSystemDatum())
			case <-stop:
				return
			}
		}
	}()

	return stop
}
