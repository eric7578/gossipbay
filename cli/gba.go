package main

import (
	"log"
	"os"

	"github.com/eric7578/gossipbay/repo"
	"github.com/eric7578/gossipbay/schedule"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gba",
		Usage: "ptt scheduled crawler",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run schedule jobs",
				Flags: []cli.Flag{
					repositoryFlag(true),
					tokenFlag(true),
					labelFlag(true),
				},
				Action: func(c *cli.Context) error {
					r := repo.NewGithub(c.String("repository"), c.String("token"))
					return schedule.Run(c.StringSlice("label")[0], r)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
