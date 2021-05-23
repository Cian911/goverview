package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	wcache "github.com/cian911/goverview/pkg/cache"
	"github.com/cian911/goverview/pkg/gh"
	"github.com/cian911/goverview/web/html"
	"github.com/google/go-github/v35/github"
	"github.com/gorilla/mux"
)

var (
	ctx          = context.Background()
	c            = gh.NewClientWithToken(ctx, os.Getenv("GITHUB_TOKEN"))
	cacheClient  = wcache.CacheClient()
	organization = "storyful"
	opts         = &github.ListWorkflowRunsOptions{ListOptions: github.ListOptions{Page: 1, PerPage: 1}}
	jobOpts      = &github.ListWorkflowJobsOptions{ListOptions: github.ListOptions{Page: 1, PerPage: 3}}
	orgOpts      = &github.RepositoryListByOrgOptions{Type: "all", Sort: "updated", Direction: "desc", ListOptions: github.ListOptions{Page: 1, PerPage: 16}}
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
	repos, _, _ := c.OrganizationRepos(ctx, organization, orgOpts)
	runs := []gh.RecentRuns{}

	for _, repo := range repos {
		run, _, _ := c.RecentWorkflowRuns(ctx, organization, *repo.Name, opts)
		if len(run.WorkflowRuns) == 0 {
			continue
		}
		recentRun := gh.RecentRuns{
			Repository: *repo.Name,
			Runs:       run.WorkflowRuns,
		}
		runs = append(runs, recentRun)
	}

	err := html.IndexPage(w, runs)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error rendering index template: %v", err))
	}

	return nil
}

func serveActions(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	key := vars["id"]
	repo := vars["repo"]
	runId, err := strconv.ParseInt(key, 10, 64)

	if err != nil {
		return NewHTTPError(err, 400, "Bad request : invalid ID.")
	}

	run, _, _ := c.WorkflowRunById(ctx, organization, repo, runId)
	jobs, _, _ := c.JobsListWorkflowRun(ctx, organization, repo, runId, jobOpts)
	data := gh.ActionData{
		Run:  run,
		Jobs: jobs,
	}

	err = html.ActionsPage(w, &data)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error rendering actions template: %v", err))
	}

	return nil
}

func HandleRoutes(router *mux.Router) {
	// Web Routes
	router.Handle("/", cacheClient.Middleware(rootHandler(serveIndex)))
	router.Handle("/workflow/{repo}/{id}", cacheClient.Middleware(rootHandler(serveActions)))

	// API routes
	router.Handle("/api/runs", rootHandler(workflowRuns))
	router.Handle("/api/runs/{id}", rootHandler(workflowRun))
	router.Handle("/api/workflows/{id}", rootHandler(workflowJob))
}

func workflowRuns(w http.ResponseWriter, r *http.Request) error {
	runs, resp, err := c.RecentWorkflowRuns(ctx, organization, "droptube-poc", opts)
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

	run, resp, err := c.WorkflowRunById(ctx, organization, "droptube-poc", runId)

	if resp.StatusCode == 404 {
		return NewHTTPError(nil, 404, "The requested workflow run was not found.")
	}

	if err != nil {
		return NewHTTPError(err, resp.StatusCode, "Error from Github API. Please check your token for the correct scopes, access rights and/or rate limits.")
	}
	json.NewEncoder(w).Encode(run)
	return nil
}

func workflowJob(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	key := vars["id"]
	runId, err := strconv.ParseInt(key, 10, 64)

	if err != nil {
		return NewHTTPError(err, 400, "Bad request : invalid ID.")
	}

	run, resp, err := c.JobsListWorkflowRun(ctx, organization, "droptube-poc", runId, jobOpts)

	if resp.StatusCode == 404 {
		return NewHTTPError(nil, 404, "The requested workflow run was not found.")
	}

	if err != nil {
		return NewHTTPError(err, resp.StatusCode, "Error from Github API. Please check your token for the correct scopes, access rights and/or rate limits.")
	}
	json.NewEncoder(w).Encode(run)
	return nil
}
