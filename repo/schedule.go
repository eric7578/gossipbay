package repo

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
	"time"

	"github.com/eric7578/gossipbay/crawler"
)

var commentTmpl *template.Template

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
	var until time.Time
	switch opt.Period {
	case TrendingWeekly:
		until = time.Now().Add(-7 * 24 * time.Hour)
	case TrendingDaily:
		until = time.Now().Add(-24 * time.Hour)
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
			c := crawler.NewCrawler(issue.Title)
			c.CollectUntil(tr, until)
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
