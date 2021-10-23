package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cian911/goverview/pkg/gh"
	"github.com/google/go-github/v35/github"
	"github.com/gorilla/websocket"
)

const tickerTime = 15 * time.Second

var (
	ctx          = context.Background()
	c            = gh.NewClientWithToken(ctx, os.Getenv("GITHUB_TOKEN"))
	opts         = &github.ListWorkflowRunsOptions{ListOptions: github.ListOptions{Page: 1, PerPage: 25}}
	jobOpts      = &github.ListWorkflowJobsOptions{ListOptions: github.ListOptions{Page: 1, PerPage: 25}}
	organization = "YOUR_ORG"
	orgOpts      = &github.RepositoryListByOrgOptions{Type: "all", Sort: "updated", Direction: "desc", ListOptions: github.ListOptions{Page: 1, PerPage: 50}}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return ws, err
	}

	return ws, nil
}

func Writer(conn *websocket.Conn, vars map[string]string) {
	jsonString := []byte{}
	path := vars["path"]

	switch path {
	case "action":
		runId, err := strconv.ParseInt(vars["id"], 10, 64)

		if err != nil {
			log.Fatalf("Could not parse websockert param: %v", err)
		}

		repository := vars["repo"]

		jsonString = actionData(repository, runId)
	case "index":
		jsonString = indexData()
	default:
		jsonString = nil
	}

	// if err != nil {
	// log.Fatalf("Failed to query data in websocket connection", err)
	// }

	for {
		ticker := time.NewTicker(tickerTime)

		for t := range ticker.C {
			fmt.Printf("Doing: %v\n", t)

			//      if *run.Status == "completed" {
			// conn.WriteMessage(websocket.TextMessage, []byte(jsonString))
			// conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			// conn.Close()
			// }

			// if err != nil {
			// fmt.Println(err)
			// }

			if err := conn.WriteMessage(websocket.TextMessage, []byte(jsonString)); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func actionData(repository string, runId int64) []byte {
	run, _, _ := c.WorkflowRunById(ctx, "YOUR_ORG", repository, runId)
	job, _, _ := c.JobsListWorkflowRun(ctx, "YOUR_ORG", repository, runId, jobOpts)
	data := gh.ActionData{
		Run:  run,
		Jobs: job,
	}

	jsonString, _ := json.Marshal(data)
	return jsonString
}

func indexData() []byte {
	// TODO: Keep list of repos stored in memory somewhere rather then making a call every time
	// TODO: Return an error here.. and above
	repos, _, _ := c.OrganizationRepos(ctx, organization, orgOpts)
	runs := []gh.RecentRuns{}

	for _, repo := range repos {
		run, _, _ := c.RecentWorkflowRuns(ctx, organization, *repo.Name, &github.ListWorkflowRunsOptions{ListOptions: github.ListOptions{Page: 1, PerPage: 1}})
		if len(run.WorkflowRuns) == 0 {
			continue
		}
		recentRun := gh.RecentRuns{
			Repository: *repo.Name,
			Runs:       run.WorkflowRuns,
		}
		runs = append(runs, recentRun)
	}

	jsonString, _ := json.Marshal(runs)
	return jsonString
}
