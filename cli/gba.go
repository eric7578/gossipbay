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
		Name:  "gba",
		Usage: "ptt scheduled crawler",
		Commands: []*cli.Command{
			cmdRun(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func cmdRun() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Run schedule jobs",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Required: true,
				Name:     "schedule",
				Aliases:  []string{"s"},
				Usage:    "schedule type",
			},
			&cli.StringFlag{
				Required: true,
				Name:     "github",
				Aliases:  []string{"gh"},
				Usage:    "github repository",
				EnvVars:  []string{"GITHUB_REPOSITORY"},
			},
			&cli.StringFlag{
				Required: true,
				Name:     "token",
				Aliases:  []string{"t"},
				Usage:    "github api token",
			},
		},
		Action: func(c *cli.Context) error {
			segs := strings.Split(c.String("github"), "/")
			r := repo.NewGithub(segs[0], segs[1], c.String("token"))
			return schedule.RunSchedule(c.String("schedule"), r)
		},
	}
}
