package main

import (
	"github.com/cloudfoundry/gosigar"
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/cloudwatch"
	"time"
)

const (
	MEMORY_UNITS = "Megabytes"
	MEMORY_UNITS_DIV = 1048576
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func GetDimensions() []cloudwatch.Dimension {
	return []cloudwatch.Dimension{
		cloudwatch.Dimension{
			Name:  "InstanceId",
			Value: aws.InstanceId(),
		},
	}
}

func GetFileSystemDatum() []cloudwatch.MetricDatum {
	dimensions := GetDimensions()

	now := time.Now()

	disk := sigar.FileSystemUsage{}
	disk.Get("/")

	return []cloudwatch.MetricDatum{
		cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "DiskSpaceUtilization",
			Unit:       "Percent",
			Timestamp:  now,
			Value:      disk.UsePercent(),
		},
		cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "DiskSpaceUsed",
			Unit:       MEMORY_UNITS,
			Timestamp:  now,
			Value:      float64(disk.Used / MEMORY_UNITS_DIV),
		},
		cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "DiskSpaceAvailable",
			Unit:       MEMORY_UNITS,
			Timestamp:  now,
			Value:      float64(disk.Free / MEMORY_UNITS_DIV),
		},
	}
}

func GetMemoryDatum() []cloudwatch.MetricDatum {
	dimensions := GetDimensions()

	now := time.Now()

	mem := sigar.Mem{}
	mem.Get()

	swap := sigar.Swap{}
	swap.Get()

	return []cloudwatch.MetricDatum{
		cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "MemoryUtilization",
			Unit:       "Percent",
			Timestamp:  now,
			Value:      float64(100 * mem.Used / mem.Total),
		},
		cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "MemoryAvailable",
			Unit:       MEMORY_UNITS,
			Timestamp:  now,
			Value:      float64(mem.Free / MEMORY_UNITS_DIV),
		},
		cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "MemoryUsed",
			Unit:       MEMORY_UNITS,
			Timestamp:  now,
			Value:      float64(mem.Used / MEMORY_UNITS_DIV),
		},
		cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "SwapUtilization",
			Unit:       "Percent",
			Timestamp:  now,
			Value:      float64(100 * swap.Used / swap.Total),
		},
		cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "SwapUsed",
			Unit:       MEMORY_UNITS,
			Timestamp:  now,
			Value:      float64(swap.Used / MEMORY_UNITS_DIV),
		},
	}
}

func main() {
	// AWS authentication
	auth, err := aws.GetAuth("", "", "", time.Now())
	check(err)

	// Initialize CloudWatch
	region := aws.Regions[aws.InstanceRegion()]
	watch, err := cloudwatch.NewCloudWatch(auth, region.CloudWatchServicepoint)
	check(err)

	// Extract metrics
	memoryDatum := GetMemoryDatum()
	fileDatum := GetFileSystemDatum()

	// Put metrics
	_, err = watch.PutMetricDataNamespace(memoryDatum, "System/Linux")
	check(err)

	_, err = watch.PutMetricDataNamespace(fileDatum, "System/Linux")
	check(err)
}
