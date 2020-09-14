package testutil

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
)

type TestDataLoader struct {
}

func (tdl *TestDataLoader) Load(p string) (*goquery.Document, error) {
	u, err := url.Parse(p)
	if err != nil {
		return nil, err
	}
	u.Scheme = ""
	u.Host = ""
	fpath := filepath.Join(MustGetwd(), u.String())
	r, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}

	return goquery.NewDocumentFromReader(r)
}

func MustGetwd() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}
