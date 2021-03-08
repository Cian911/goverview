package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/cian911/goverview/pkg/gh"
	"github.com/cian911/goverview/web/html"
	"github.com/google/go-github/v33/github"
	"github.com/gorilla/mux"
)

var (
	ctx  = context.Background()
	c    = gh.NewClientWithToken(ctx, os.Getenv("GITHUB_TOKEN"))
	opts = &github.ListWorkflowRunsOptions{ListOptions: github.ListOptions{Page: 1, PerPage: 25}}
)

type rootHandler func(http.ResponseWriter, *http.Request) error

func (fn rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Call handler function
	err := fn(w, r)
	if err == nil {
		return
	}

	// An error occured, start logging process
	log.Printf("An error occured: %v", err)

	// Check if it's a client error
	clientError, ok := err.(ClientError)
	if !ok {
		// If not a client error, assume it's a server error..
		w.WriteHeader(500)
		return
	}

	// Try and get ResponseBody from clientError
	body, err := clientError.ResponseBody()
	if err != nil {
		log.Printf("An error occured: %v", err)
		w.WriteHeader(500)
		return
	}

	// Get http status code and headers
	status, headers := clientError.ResponseHeaders()
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
	w.Write(body)
}

func serveIndex(w http.ResponseWriter, r *http.Request) error {
	// lp := filepath.Join("./web/html/layouts", "layout.html")
	// fp := filepath.Join("./web/html", "index.html")
	//
	// tpl, err := template.ParseFiles(lp, fp)
	// if err != nil {
	//   fmt.Errorf("%v", err)
	//   return err
	// }
	runs, _, _ := c.RecentWorkflowRuns(ctx, "Cian911", "gomerge", opts)
	html.IndexPage(w, runs)
	return nil
}

func HandleRoutes(router *mux.Router) {
	// Web Routes
	router.Handle("/", rootHandler(serveIndex))
	// router.Handle("/workflow/{id}", rootHandler(serveWorkflow))

	// API routes
	router.Handle("/api/runs", rootHandler(workflowRuns))
	router.Handle("/api/runs/{id}", rootHandler(workflowRun))
}

func workflowRuns(w http.ResponseWriter, r *http.Request) error {
	runs, resp, err := c.RecentWorkflowRuns(ctx, "Cian911", "gomerge", opts)
	if err != nil {
		return NewHTTPError(err, resp.StatusCode, "Error from Github API. Please check your token for the correct scopes, access rights and/or rate limits.")
	}

	json.NewEncoder(w).Encode(runs)
	return nil
}

func workflowRun(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	key := vars["id"]
	runId, err := strconv.ParseInt(key, 10, 64)

	if err != nil {
		return NewHTTPError(err, 400, "Bad request : invalid ID.")
	}

	run, resp, err := c.WorkflowRunById(ctx, "Cian911", "gomerge", runId)

	if resp.StatusCode == 404 {
		return NewHTTPError(nil, 404, "The requested workflow run was not found.")
	}

	if err != nil {
		return NewHTTPError(err, resp.StatusCode, "Error from Github API. Please check your token for the correct scopes, access rights and/or rate limits.")
	}
	json.NewEncoder(w).Encode(run)
	return nil
}
