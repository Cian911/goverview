package html

import (
	"embed"
	"html/template"
	"io"

	"github.com/google/go-github/v33/github"
)

//go:embed *
var files embed.FS

func IndexPage(w io.Writer, runs *github.WorkflowRuns) error {
	index := parse("index.html")
	return index.Execute(w, runs)
}

func ActionsPage(w io.Writer, jobs *github.Jobs) error {
	index := parse("actions.html")
	return index.Execute(w, jobs)
}
func parse(file string) *template.Template {
	return template.Must(template.New("layout.html").ParseFS(files, "layout.html", file))
}
