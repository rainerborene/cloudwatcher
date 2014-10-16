package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertError(t *testing.T, msg string) {
	if err := Config.Valid(); err != nil {
		assert.Equal(t, err.Error(), msg)
	}
	Config.DiskPath = "/"
	Config.InstanceId = "i3190"
	Config.MemoryUnits = "Megabytes"
	Config.DiskSpaceUnits = "Gigabytes"
}

func TestValid(t *testing.T) {
	// Validate presence of instance id
	Config.InstanceId = ""
	assertError(t, "Cannot obtain instance id from EC2 meta-data.")

	// Validate presence of memory units
	Config.MemoryUnits = ""
	assertError(t, "Value of memory units is not specified.")

	// Validate disk space units
	Config.DiskSpaceUnits = ""
	assertError(t, "Value of disk space units is not specified.")

	// Validate inclusion of memory units
	Config.MemoryUnits = "zz"
	assertError(t, "Unsupported memory units 'zz'. Use Bytes, Kilobytes, Megabytes, or Gigabytes.")

	// Validate inclusion of disk space units
	Config.DiskSpaceUnits = "zz"
	assertError(t, "Unsupported disk space units 'zz'. Use Bytes, Kilobytes, Megabytes, or Gigabytes.")

	// Validate presence of disk path
	Config.DiskPath = ""
	assertError(t, "Value of disk path is not specified.")

	// Validate existence of disk path
	Config.DiskPath = "/z"
	assertError(t, "Disk file path '/z' does not exist or cannot be accessed.")

	// Check disk and memory functions
	Config.MemAvail = false
	Config.MemUsed = false
	Config.MemUtil = false
	assertError(t, "No metrics specified for collection and submission to CloudWatch.")
}
