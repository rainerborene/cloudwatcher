package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.1.0"
	app.Name = "cloudwatcher"
	app.Usage = "Reports memory, swap, and disk space utilization metrics for an Amazon EC2 Linux instance."
	app.Commands = []cli.Command{
		{
			Name: "stats",
			Usage: "Displays the most recent utilization statistics",
			ShortName: "s",
			Action: func(c *cli.Context) {
				return
			},
		},
		{
			Name: "watch",
			Usage: "Collects system metrics on an Amazon EC2 instance",
			ShortName: "w",
			Action: func(c *cli.Context) {
				collector := NewCollector(Config.Interval)
				<-collector.Run()
			},
		},
	}

	app.Run(os.Args)
}
