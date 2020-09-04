package schedule

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/eric7578/gossipbay/crawler"
)

func (s *Scheduler) Run(opt RunOption) (BoardReport, error) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		report BoardReport
	)
	if opt.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), opt.Timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	posts := make([]crawler.Post, 0)
	postc, errc := s.run(ctx, opt)
	for {
		select {
		case p, ok := <-postc:
			posts = append(posts, p)
			if !ok {
				goto END
			}
		case err := <-errc:
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					goto END
				}
				return report, err
			}
		}
	}

END:
	tr := crawler.NewTrending(posts)
	report = BoardReport{
		RunOption: opt,
		Total:     len(posts),
		Threads:   tr.Deviate(opt.Deviate),
	}
	return report, nil
}

func (s *Scheduler) run(ctx context.Context, opt RunOption) (<-chan crawler.Post, <-chan error) {
	var (
		page  = fmt.Sprintf("/bbs/%s/index.html", opt.Board)
		postc = make(chan crawler.Post)
		errc  = make(chan error)
	)

	go func() {
		for page != "" {
			select {
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					errc <- err
				}
				goto END

			default:
				var (
					wg    sync.WaitGroup
					infos []crawler.PostInfo
				)
				infos, page, _ = s.crawler.VisitBoard(page, opt.From, opt.To)
				wg.Add(len(infos))
				for _, info := range infos {
					go func(info crawler.PostInfo) {
						defer wg.Done()
						if p, err := s.crawler.VisitPost(info.URL); err != nil {
							errc <- err
						} else {
							postc <- p
						}
					}(info)
				}
				wg.Wait()
			}
		}

	END:
		close(postc)
		close(errc)
	}()

	return postc, errc
}
