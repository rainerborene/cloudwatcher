package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	err := Config.Parse()
	check(err)

	app := cli.NewApp()
	app.Version = "0.2.0"
	app.Name = "cloudwatcher"
	app.Usage = "Reports memory, swap, and disk space utilization metrics for an Amazon EC2 Linux instance."
	app.Commands = []cli.Command{
		{
			Name:      "statistics",
			Usage:     "Displays the most recent utilization statistics",
			ShortName: "s",
			Action: func(c *cli.Context) {
				return
			},
		},
	}

	app.Action = func(c *cli.Context) {
		collector := NewCollector(Config.Interval)
		<-collector.Run()
	}

	app.Run(os.Args)
}
