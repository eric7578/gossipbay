package main

import (
	"log"
	"os"

	"github.com/eric7578/gossipbay/flag"
	"github.com/eric7578/gossipbay/repo"
	"github.com/eric7578/gossipbay/schedule"
	cli "github.com/urfave/cli/v2"
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
					flag.NewRepositoryFlag(true),
					flag.NewTokenFlag(true),
					flag.NewLabelFlag(true),
				},
				Action: func(c *cli.Context) error {
					r := repo.NewGithub(c.String("repository"), c.String("token"))
					return schedule.Run(r, c.StringSlice("label")[0])
				},
			},
			{
				Name:  "prune",
				Usage: "Remove obsoleted comments",
				Flags: []cli.Flag{
					flag.NewRepositoryFlag(true),
					flag.NewTokenFlag(true),
					&cli.StringFlag{
						Name:        "range",
						Usage:       "Comments created `DAYS` days ago",
						Required:    true,
						DefaultText: ":",
					},
					&cli.StringFlag{
						Name:     "user",
						Usage:    "`USER` who create the comment",
						EnvVars:  []string{"GITHUB_ACTOR"},
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					r := repo.NewGithub(c.String("repository"), c.String("token"))
					from, to := flag.ParseDaysExpression(c.String("range"))
					return schedule.Prune(r, c.String("user"), from, to)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
