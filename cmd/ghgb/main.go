package main

import (
	"os"
	"strings"

	"github.com/eric7578/gossipbay/repo"
	"github.com/eric7578/gossipbay/schedule"
)

func main() {
	token := os.Getenv("GB_TOKEN")
	segs := strings.Split(os.Getenv("GB_REPOSITORY"), "/")
	r := repo.NewGithub(segs[0], segs[1], token)
	schedule.RunSchedule("trending-daily", r)
}
