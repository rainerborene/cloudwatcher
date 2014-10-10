package main

import (
	"github.com/cloudfoundry/gosigar"
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/cloudwatch"
	"time"
)

const MEMORY_UNITS_DIV = 1048576

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func GetMemoryDatum() []cloudwatch.MetricDatum {
	now := time.Now()
	mem := sigar.Mem{}
	mem.Get()

	dimensions := []cloudwatch.Dimension{
		cloudwatch.Dimension{
			Name:  "InstanceId",
			Value: aws.InstanceId(),
		},
	}

	return []cloudwatch.MetricDatum{
		cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "MemoryUtilization",
			Unit:       "Megabytes",
			Timestamp:  now,
			Value:      float64(100 * mem.Used / mem.Total),
		},
		cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "MemoryAvailable",
			Unit:       "Megabytes",
			Timestamp:  now,
			Value:      float64(mem.Free / MEMORY_UNITS_DIV),
		},
		cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "MemoryUsed",
			Unit:       "Megabytes",
			Timestamp:  now,
			Value:      float64(mem.Used / MEMORY_UNITS_DIV),
		},
	}
}

func main() {
	// AWS authentication
	auth, err := aws.GetAuth("", "", "", time.Now())
	check(err)

	region := aws.Regions[aws.InstanceRegion()]

	// // Initialize CloudWatch
	watch, err := cloudwatch.NewCloudWatch(auth, region.CloudWatchServicepoint)
	check(err)

	memoryDatum := GetMemoryDatum()

	_, err = watch.PutMetricData(memoryDatum)
	check(err)
}
