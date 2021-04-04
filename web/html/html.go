package html

import (
	"embed"
	"html/template"
	"io"
	"time"

	"github.com/cian911/goverview/pkg/gh"
	"github.com/google/go-github/v33/github"
	"github.com/xeonx/timeago"
)

//go:embed *
var files embed.FS

func IndexPage(w io.Writer, runs *github.WorkflowRuns) error {
	index := parse("index.html")
	return index.Execute(w, runs)
}

func ActionsPage(w io.Writer, jobs *gh.ActionData) error {
	index := parse("actions.html")
	return index.Execute(w, jobs)
}

func parse(file string) *template.Template {
	funcs := template.FuncMap{
		"toString": func(str *string) string {
			return *str
		},
		"timeAgo": func(t *github.Timestamp) string {
			return timeago.NoMax(timeago.English).FormatReference(t.Time, time.Now())
		},
	}

	return template.Must(template.New("layout.html").Funcs(funcs).ParseFS(files, "layout.html", file))
}
