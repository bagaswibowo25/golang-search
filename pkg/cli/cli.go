package cli

import (
	"log"
	"os"

	"github.com/bagaswibowo25/golang-search/pkg/app"
	"github.com/bagaswibowo25/golang-search/pkg/config"
	cli "github.com/urfave/cli/v2"
)

func Run() {
	run := newSearch()
	if err := run.Run(os.Args); err != nil {
		log.Printf("Error running command: %s", err)
	}
}

func newSearch() *cli.App {
	search := &cli.App{
		Name:  "golang-search",
		Usage: "CLI tool",
		Commands: []*cli.Command{
			{
				Name:  "start",
				Usage: "Start golang-search server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "Port to listen",
						Value:   "8080",
					},
					&cli.StringFlag{
						Name:     "consumer-id",
						Aliases:  []string{"c"},
						Usage:    "c",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "nats-url",
						Aliases: []string{"ns"},
						Usage:   "Nats URL",
						Value:   "https://localhost:4222",
					},
					&cli.StringFlag{
						Name:     "subject",
						Usage:    "NATS Subject",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "stream",
						Usage:    "NATS Stream",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					httpConf := &config.HttpConfig{
						ListenPort: ":" + c.String("port"),
					}

					natsConf := &config.NatsConfig{
						ConsumerId: c.String("consumer-id"),
						NatsURL:    c.String("nats-url"),
						Subject:    c.String("subject"),
						Stream:     c.String("stream"),
					}

					err := app.StartServer(httpConf, natsConf)
					if err != nil {
						log.Printf("cant start golang-search %s", err)
					}

					return err
				},
			},
		},
	}
	return search
}
