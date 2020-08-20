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
		Name:  "gossipbay",
		Usage: "ptt scheduled crawler",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Required: true,
				Name:     "owner",
				Aliases:  []string{"o"},
				Usage:    "repository owner",
				EnvVars:  []string{"GH_OWNER"},
			},
			&cli.StringFlag{
				Required: true,
				Name:     "repository",
				Aliases:  []string{"r"},
				Usage:    "repository name",
				EnvVars:  []string{"GB_REPO"},
			},
			&cli.StringFlag{
				Required: true,
				Name:     "token",
				Aliases:  []string{"t"},
				Usage:    "github api token",
				EnvVars:  []string{"GH_TOKEN"},
			},
			&cli.BoolFlag{
				Name:    "github",
				Aliases: []string{"gh"},
				Usage:   "use github",
			},
		},
		Action: func(c *cli.Context) error {
			owner, repository, token := c.String("owner"), c.String("repository"), c.String("token")

			var r repo.Repository
			switch {
			case c.Bool("github"):
				r = repo.NewGithub(owner, repository, token)
			}

			return schedule.RunSchedule("trending-daily", r)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
