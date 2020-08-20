package main

import (
	"log"
	"os"
	"strings"

	"github.com/eric7578/gossipbay/repo"
	"github.com/eric7578/gossipbay/schedule"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gossipbay",
		Usage: "ptt scheduled crawler",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "github",
				Aliases: []string{"gh"},
				Usage:   "github repository",
				EnvVars: []string{"GITHUB_REPOSITORY"},
			},
			&cli.StringFlag{
				Required: true,
				Name:     "token",
				Aliases:  []string{"t"},
				Usage:    "github api token",
			},
			&cli.StringFlag{
				Required: true,
				Name:     "schedule",
				Aliases:  []string{"s"},
				Usage:    "schedule type",
			},
		},
		Action: func(c *cli.Context) error {
			segs := strings.Split(c.String("github"), "/")
			r := repo.NewGithub(segs[0], segs[1], c.String("token"))
			return schedule.RunSchedule(c.String("schedule"), r)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
