package main

import (
	"log"
	"os"

	"github.com/eric7578/gossipbay/flagutil"
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
					newRepositoryFlag(true),
					newTokenFlag(true),
					newLabelFlag(true),
				},
				Action: func(c *cli.Context) error {
					r := repo.NewGithub(c.String("repository"), c.String("token"))
					s := c.StringSlice("label")[0]
					from, to := flagutil.ParseSchedule(s)
					return schedule.Run(r, s, from, to)
				},
			},
			{
				Name:  "prune",
				Usage: "Remove obsoleted comments",
				Flags: []cli.Flag{
					newRepositoryFlag(true),
					newTokenFlag(true),
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
					from, to := flagutil.ParseDaysExpression(c.String("range"))
					return schedule.Prune(r, c.String("user"), from, to)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func newRepositoryFlag(required bool) *cli.StringFlag {
	return &cli.StringFlag{
		Required: required,
		Name:     "repository",
		Aliases:  []string{"r"},
		Usage:    "Repository OWNER/REPO",
		EnvVars:  []string{"GITHUB_REPOSITORY"},
	}
}

func newTokenFlag(required bool) *cli.StringFlag {
	return &cli.StringFlag{
		Required: required,
		Name:     "token",
		Aliases:  []string{"t"},
		Usage:    "Github api token",
	}
}

func newLabelFlag(required bool) *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Required: required,
		Name:     "label",
		Aliases:  []string{"l"},
		Usage:    "Issue flags",
	}
}
