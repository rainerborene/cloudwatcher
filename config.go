package main

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/crowdmob/goamz/aws"
	"os"
	"time"
)

var paths = []string{
	"/etc/cloudwatcher.toml",
	"cloudwatcher.toml",
}

var units = map[string]uint64{
	"Bytes":     1,
	"Kilobytes": 1024,
	"Megabytes": 1048576,
	"Gigabytes": 1073741824,
}

type config struct {
	Namespace      string        `toml:"namespace"`
	InstanceId     string        `toml:"instance_id"`
	MemUtil        bool          `toml:"mem_util"`
	MemUsed        bool          `toml:"mem_used"`
	MemAvail       bool          `toml:"mem_avail"`
	SwapUtil       bool          `toml:"swap_util"`
	SwapUsed       bool          `toml:"swap_used"`
	DiskPath       string        `toml:"disk_path"`
	DiskSpaceUtil  bool          `toml:"disk_space_util"`
	DiskSpaceUsed  bool          `toml:"disk_space_used"`
	DiskSpaceAvail bool          `toml:"disk_space_avail"`
	MemoryUnits    string        `toml:"memory_units"`
	DiskSpaceUnits string        `toml:"disk_space_units"`
	Interval       time.Duration `toml:"interval"`
}

func (c *config) MemoryEnabled() bool {
	return c.MemAvail || c.MemUsed || c.MemUtil
}

func (c *config) DiskEnabled() bool {
	return c.DiskSpaceAvail || c.DiskSpaceUsed || c.DiskSpaceUtil
}

func (c *config) DiskSpaceUnitsInt() uint64 {
	return units[c.DiskSpaceUnits]
}

func (c *config) MemoryUnitsInt() uint64 {
	return units[c.MemoryUnits]
}

func (c *config) Valid() error {
	if c.InstanceId == "" {
		return errors.New("Cannot obtain instance id from EC2 meta-data.")
	} else if c.MemoryUnits == "" {
		return errors.New("Value of memory units is not specified.")
	} else if c.DiskSpaceUnits == "" {
		return errors.New("Value of disk space units is not specified.")
	} else if c.MemoryUnitsInt() == 0 {
		return errors.New(fmt.Sprintf("Unsupported memory units '%s'. Use Bytes, Kilobytes, Megabytes, or Gigabytes.", c.MemoryUnits))
	} else if c.DiskSpaceUnitsInt() == 0 {
		return errors.New(fmt.Sprintf("Unsupported disk space units '%s'. Use Bytes, Kilobytes, Megabytes, or Gigabytes.", c.DiskSpaceUnits))
	} else if c.DiskPath == "" && c.DiskEnabled() {
		return errors.New("Value of disk path is not specified.")
	} else if _, err := os.Stat(c.DiskPath); err != nil && c.DiskEnabled() {
		return errors.New(fmt.Sprintf("Disk file path '%s' does not exist or cannot be accessed.", c.DiskPath))
	} else if !c.DiskEnabled() && !c.MemoryEnabled() {
		return errors.New("No metrics specified for collection and submission to CloudWatch.")
	}

	return nil
}

func (c *config) Path() string {
	for i := range paths {
		if _, err := os.Stat(paths[i]); err == nil {
			return paths[i]
		}
	}
	return ""
}

func (c *config) Parse() error {
	if path := c.Path(); path != "" {
		if _, err := toml.DecodeFile(path, c); err != nil {
			return err
		}
	}

	if c.InstanceId == "" {
		c.InstanceId = aws.InstanceId()
	}

	return c.Valid()
}

var Config = &config{
	MemUtil:        true,
	MemUsed:        true,
	MemAvail:       true,
	SwapUtil:       true,
	SwapUsed:       true,
	DiskPath:       "/",
	DiskSpaceUtil:  true,
	DiskSpaceUsed:  true,
	DiskSpaceAvail: true,
	MemoryUnits:    "Megabytes",
	DiskSpaceUnits: "Gigabytes",
	Namespace:      "System/Linux",
	Interval:       60,
}
