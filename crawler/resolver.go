package crawler

import (
	"fmt"
	"regexp"
)

var (
	regProtocol = regexp.MustCompile("https?://")
)

type resolver struct {
	domain string
}

func (r *resolver) getBoardIndex(board string) string {
	return fmt.Sprintf("%s/bbs/%s/index.html", r.domain, board)
}

func (r *resolver) getFullURL(path string) string {
	if regProtocol.MatchString(path) {
		return path
	}
	return r.domain + path
}
