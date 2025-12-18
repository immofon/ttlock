package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/immofon/ttlock"
	"github.com/urfave/cli/v2"
)

var client *ttlock.Client

func main() {
	app := &cli.App{
		Name:  "ttlock",
		Usage: "TTLock CLI",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Config toml file path",
				Value:   "/tmp/ttlock.toml",
			},
		},
		Before: func(ctx *cli.Context) error {
			config_path := ctx.String("config")
			var config Config

			if _, err := os.Stat(config_path); err != nil {
				log.Printf("failed to load config file: %v", err)
				data, err := toml.Marshal(config)
				if err != nil {
					return err
				}

				return os.WriteFile(config_path, data, 0600)
			}
			if _, err := toml.DecodeFile(config_path, &config); err != nil {
				return err
			}

			client = ttlock.NewClient(
				config.ClientID,
				config.ClientSecret,
				config.Username,
				config.Password,
			)
			return nil
		},
		After: func(ctx *cli.Context) error {
			return nil
		},
		Commands: []*cli.Command{
			helloCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
