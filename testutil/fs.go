package testutil

import (
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
)

type TestDataLoader struct {
}

func (tdl *TestDataLoader) Load(p string) (*goquery.Document, error) {
	fpath := filepath.Join(MustGetwd(), p)
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
