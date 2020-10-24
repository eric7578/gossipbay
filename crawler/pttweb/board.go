package pttweb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type scanResult struct {
	Post Post
	Err  error
}

func (w *PttWorker) loadBoardPage(page string, from time.Time, to time.Time) ([]string, string, error) {
	doc, err := w.Load(getFullURL(page))
	posts := make([]string, 0)
	if err != nil {
		return nil, "", err
	}

	next := true
	doc.
		Find(".r-list-container").
		Children().
		Filter(".search-bar").
		NextUntil(".r-list-sep").
		Each(func(i int, sel *goquery.Selection) {
			title := sel.Find(".title > a")
			href, ok := title.Attr("href")
			if !ok {
				return
			}

			_, createdAt := parseURL(href)
			if createdAt.Before(from) {
				next = false
				return
			} else if createdAt.Before(to) {
				posts = append(posts, getFullURL(href))
			}
		})

	if next {
		nextHref := doc.Find(".btn-group-paging .btn").Eq(1).AttrOr("href", "")
		return posts, getFullURL(nextHref), nil
	}

	return posts, "", nil
}

func (w *PttWorker) scanBoard(ctx context.Context, board string, from, to time.Time) <-chan scanResult {
	boardPage := fmt.Sprintf("/bbs/%s/index.html", board)
	resultc := make(chan scanResult)
	var wg sync.WaitGroup

	go func() {
		for boardPage != "" {
			select {
			case <-ctx.Done():
				goto END

			default:
				var pages []string
				pages, boardPage, _ = w.loadBoardPage(boardPage, from, to)
				for _, page := range pages {
					wg.Add(1)
					go func(page string) {
						defer wg.Done()
						p, err := w.VisitPost(VisitPostOption{
							URL: page,
						})
						resultc <- scanResult{
							Post: p,
							Err:  err,
						}
					}(page)
				}
			}
		}

	END:
		wg.Wait()
		close(resultc)
	}()

	return resultc
}