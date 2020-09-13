package main

import (
	"log"
	"os"
	"time"

	"github.com/eric7578/gossipbay/crawler"
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
				Name:  "visit-post",
				Usage: "Run crawler on an single page by url",
				Action: func(c *cli.Context) error {
					pageURL := c.Args().First()
					cr := crawler.NewPageCrawler()
					post, err := cr.VisitPost(pageURL)
					if err != nil {
						return err
					}

					return schedule.Pipe(post, os.Stdout)
				},
			},
			{
				Name:  "trending",
				Usage: "Run a single board job",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "board",
						Aliases:  []string{"b"},
						Usage:    "Name of the target board",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "schedule",
						Aliases:  []string{"s"},
						Usage:    "Job scuedule",
						Required: true,
					},
					&cli.Float64Flag{
						Name:     "deviate",
						Aliases:  []string{"d"},
						Usage:    "Deviate for calculate trending",
						Required: true,
					},
					&cli.Int64Flag{
						Name:    "timeout",
						Aliases: []string{"t"},
						Usage:   "Timeout for crawling",
					},
				},
				Action: func(c *cli.Context) (err error) {
					from, to, err := flagutil.ParseSchedule(c.String("schedule"))
					if err != nil {
						return err
					}

					opt := schedule.TrendingOption{
						Board:   c.String("board"),
						From:    from,
						To:      to,
						Timeout: time.Second * time.Duration(c.Int64("timeout")),
						Deviate: c.Float64("deviate"),
					}
					s := schedule.NewScheduler()
					threads, err := s.Trending(opt)
					if err != nil {
						return err
					}

					return schedule.Pipe(threads, os.Stdout)
				},
			},
			{
				Name:  "repo",
				Usage: "Run crawler based on repoistory settings",
				Subcommands: []*cli.Command{
					{
						Name:  "run",
						Usage: "Run schedule jobs",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Required: true,
								Name:     "repository",
								Aliases:  []string{"r"},
								Usage:    "Repository OWNER/REPO",
								EnvVars:  []string{"GITHUB_REPOSITORY"},
							},
							&cli.StringFlag{
								Name:    "token",
								Aliases: []string{"t"},
								Usage:   "Github api token",
							},
							&cli.StringSliceFlag{
								Required: true,
								Name:     "label",
								Aliases:  []string{"l"},
								Usage:    "Issue flags",
							},
						},
						Action: func(c *cli.Context) error {
							r := repo.NewGithub(c.String("repository"), c.String("token"))
							s := schedule.NewScheduler()
							report, err := s.RunIssues(r, schedule.RunIssueOptions{
								Labels: c.StringSlice("label"),
							})
							if err != nil {
								return err
							}

							return schedule.Pipe(report, os.Stdout)
						},
					},
					{
						Name: "prune",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Required: true,
								Name:     "repository",
								Aliases:  []string{"r"},
								Usage:    "Repository OWNER/REPO",
								EnvVars:  []string{"GITHUB_REPOSITORY"},
							},
							&cli.StringFlag{
								// Required: true,
								Name:    "token",
								Aliases: []string{"t"},
								Usage:   "Github api token",
							},
							&cli.IntFlag{
								Required: true,
								Name:     "days-ago",
								Usage:    "Prune artifact `DAYS` days ago",
							},
						},
						Action: func(c *cli.Context) error {
							r := repo.NewGithub(c.String("repository"), c.String("token"))
							return r.PruneArtifact(c.Int("days-ago"))
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
