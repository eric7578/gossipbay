package repo

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
	"time"

	"github.com/eric7578/gossipbay/crawler"
)

var (
	commentTmpl *template.Template
	taipei      *time.Location
)

func init() {
	comment := `{{ range . }}
{{ range .Posts }}[{{ .Title }}]({{ .URL }}) **{{ .NumUp }}** 推 **{{ .NumDown }}** 噓
{{end}}
{{end}}
`
	var err error
	commentTmpl, err = template.New("comment").Parse(comment)
	if err != nil {
		panic(err)
	}

	taipei, err = time.LoadLocation("Asia/Taipei")
	if err != nil {
		panic(err)
	}
}

type schedule struct {
	Repository
	until   time.Time
	deviate float32
}

type ScheduleOption struct {
	Period  string
	Deviate float32
}

func RunSchedule(r Repository, opt ScheduleOption) {
	now := time.Now()
	var (
		from time.Time
		to   time.Time = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, taipei)
	)
	switch opt.Period {
	case TrendingWeekly:
		from = to.Add(-7 * 24 * time.Hour)
	case TrendingDaily:
		from = to.Add(-24 * time.Hour)
	default:
		panic(fmt.Errorf("invalid period %s", opt.Period))
	}

	issues := r.ListIssues(opt.Period)
	var wg sync.WaitGroup
	wg.Add(len(issues))
	for _, issue := range issues {
		go func(issue Issue) {
			defer wg.Done()
			tr := crawler.NewTrending()
			c := crawler.NewCrawler()
			c.Collect(tr, crawler.CollectOption{
				Board: issue.Title,
				From:  from,
				To:    to,
			})
			r.CreateIssueComment(issue.ID, generateComment(tr.Deviate(opt.Deviate)))
		}(issue)
	}
	wg.Wait()
}

func generateComment(threads []crawler.Thread) string {
	var tpl bytes.Buffer
	err := commentTmpl.Execute(&tpl, threads)
	if err != nil {
		panic(err)
	}
	return tpl.String()
}
