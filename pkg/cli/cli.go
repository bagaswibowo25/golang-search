package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/bagaswibowo25/golang-search/pkg/app"
	"github.com/bagaswibowo25/golang-search/pkg/config"
	cli "github.com/urfave/cli/v2"
)

func Run() {
	cmd := startCmd()
	if err := cmd.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func startCmd() *cli.App {
	err := &cli.App{
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
					&cli.StringFlag{
						Name:    "nats-url",
						Aliases: []string{"ns"},
						Usage:   "Nats URL",
					},
					&cli.StringFlag{
						Name:  "subject",
						Usage: "Nats URL",
					},
					&cli.StringFlag{
						Name:  "stream",
						Usage: "Nats URL",
					},
				},
				Action: func(c *cli.Context) error {
					httpConf := &config.HttpConfig{
						ListenPort: ":" + c.String("port"),
					}

					natsConf := &config.NatsConfig{
						ConsumerId: ":" + c.String("consumer-id"),
						NatsURL:    c.String("nats-url"),
						Subject:    c.String("subject"),
						Stream:     c.String("stream"),
					}

					if natsConf.NatsURL == "" {
						natsConf.NatsURL = "http://localhost:4222"
					}

					if natsConf.Subject == "" {
						natsConf.Subject = "loggerSubject"
					}

					if natsConf.Stream == "" {
						natsConf.Stream = "loggerStream"
					}

					start := app.StartServer(httpConf, natsConf)

					if start != nil {
						log.Fatalf("cant start golang-search %s", start)
					}

					return nil
				},
			},
		},
	}
	return err
}
