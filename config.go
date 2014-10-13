package main

import (
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

func (c *config) DiskSpaceUnitsInt() uint64 {
	return units[c.DiskSpaceUnits]
}

func (c *config) MemoryUnitsInt() uint64 {
	return units[c.MemoryUnits]
}

func (c *config) Valid() error {
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
	Interval:       time.Second,
}
