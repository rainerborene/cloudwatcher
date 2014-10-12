package main

import (
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/cloudwatch"
	"log"
	"time"
)

type collector struct {
	region   aws.Region
	server   *cloudwatch.CloudWatch
	duration time.Duration
	rotation <-chan time.Time
}

func NewCollector(duration time.Duration) *collector {
	region := aws.Regions[aws.InstanceRegion()]
	instance := &collector{duration: duration, region: region}
	instance.RotateCredentials()
	return instance
}

func (c *collector) RotateCredentials() {
	// IAM Authentication
	auth, err := aws.GetAuth("", "", "", time.Now())
	check(err)

	// Initialize CloudWatch
	server, err := cloudwatch.NewCloudWatch(auth, c.region.CloudWatchServicepoint)
	check(err)

	// Security credentials are temporary and EC2 rotate them automatically.
	c.rotation = time.After(time.Since(auth.Expiration().Add(-time.Minute)))
	c.server = server
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
			case <-c.rotation:
				c.RotateCredentials()
			case <-stop:
				return
			}
		}
	}()

	return stop
}
