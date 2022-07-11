package main

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"

	"github.com/diwufeiwen/venus-metrics/version"
)

var log = logging.Logger("main")

func main() {
	_ = logging.SetLogLevel("*", "INFO")

	app := &cli.App{
		Name:                 "venus-metrics",
		Usage:                "for venus chain service metrics",
		Version:              version.UserVersion,
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "listen",
				Usage: "service listening address(ip:port)",
				Value: "0.0.0.0:4567",
			},
		},

		Commands: []*cli.Command{
			runCmd,
		},
	}
	app.Setup()

	RunApp(app)
}
