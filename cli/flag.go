package main

import "github.com/urfave/cli/v2"

func repositoryFlag(required bool) *cli.StringFlag {
	return &cli.StringFlag{
		Required: required,
		Name:     "repository",
		Aliases:  []string{"r"},
		Usage:    "Repository OWNER/REPO",
		EnvVars:  []string{"GITHUB_REPOSITORY"},
	}
}

func tokenFlag(required bool) *cli.StringFlag {
	return &cli.StringFlag{
		Required: required,
		Name:     "token",
		Aliases:  []string{"t"},
		Usage:    "Github api token",
	}
}

func labelFlag(required bool) *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Required: required,
		Name:     "label",
		Aliases:  []string{"l"},
		Usage:    "Issue flags",
	}
}
