package html

import (
	"embed"
	"html/template"
	"io"
	"strings"
	"time"

	"github.com/cian911/goverview/pkg/gh"
	"github.com/google/go-github/v34/github"
	"github.com/xeonx/timeago"
)

//go:embed *
var files embed.FS

func IndexPage(w io.Writer, runs []gh.RecentRuns) error {
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
		"compact": func(str *string) string {
			str1 := strings.Replace(*str, " ", "-", -1)
			str2 := strings.Replace(str1, "@", "", -1)
			str3 := strings.Replace(str2, "/", "", -1)
			str4 := strings.Replace(str3, ":", "", -1)
			return strings.ToLower(str4)
		},
		"colorCase": func(str *string) string {
			switch *str {
			case "cancelled", "failure", "timed_out", "startup_failure":
				return "bg-dang"
			case "queued", "in_progress", "waiting", "requested", "skipped":
				return "bg-warn"
			case "success":
				return "bg-succ"
			default:
				return "bg-warn"
			}
		},
		"iconCase": func(str *string) string {
			switch *str {
			case "cancelled", "failure", "timed_out", "startup_failure":
				return "bi-x"
			case "queued", "in_progress":
				return "bi-arrow-clockwise icn-spinner"
			case "waiting", "requested", "skipped":
				return "bi-exclamation-circle"
			case "success":
				return "bi-check"
			default:
				return "bi-exclamation-circle"
			}
		},
	}

	return template.Must(template.New("layout.html").Funcs(funcs).ParseFS(files, "layout.html", file))
}
