package cli

import (
	"fmt"
	"os"

	"github.com/bagaswibowo25/golang-search/pkg/app"
	"github.com/bagaswibowo25/golang-search/pkg/config"
	cli "github.com/urfave/cli/v2"
)

func Run() {
	app := &cli.App{
		Name:  "golang-search",
		Usage: "CLI tool",
		Commands: []*cli.Command{
			{
				Name:  "start",
				Usage: "Start golang-search server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "port",
						Aliases:  []string{"p"},
						Usage:    "Port to listen",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "consumer-id",
						Aliases:  []string{"c"},
						Usage:    "c",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					conf := &config.Config{
						ListenPort: ":" + c.String("port"),
						ConsumerID: c.String("consumer-id"),
					}

					start := app.StartServer(conf.ListenPort, conf.ConsumerID)

					if start != nil {
						fmt.Errorf("cant start golang-search %s", start)
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
