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
	"github.com/google/go-github/v33/github"
	"github.com/gorilla/websocket"
)

const tickerTime = 15 * time.Second

var (
	ctx     = context.Background()
	c       = gh.NewClientWithToken(ctx, os.Getenv("GITHUB_TOKEN"))
	opts    = &github.ListWorkflowRunsOptions{ListOptions: github.ListOptions{Page: 1, PerPage: 25}}
	jobOpts = &github.ListWorkflowJobsOptions{ListOptions: github.ListOptions{Page: 1, PerPage: 25}}
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
	runId, err := strconv.ParseInt(vars["id"], 10, 64)

	if err != nil {
		log.Fatalf("Could not parse websockert param: %v", err)
	}

	repository := vars["repo"]

	for {
		ticker := time.NewTicker(tickerTime)

		for t := range ticker.C {
			fmt.Printf("Doing: %v\n", t)

			run, _, _ := c.WorkflowRunById(ctx, "storyful", repository, runId)
			jobs, _, _ := c.JobsListWorkflowRun(ctx, "storyful", repository, runId, jobOpts)
			data := gh.ActionData{
				Run:  run,
				Jobs: jobs,
			}

			jsonString, err := json.Marshal(data)
			if err != nil {
				fmt.Println(err)
			}

			if err := conn.WriteMessage(websocket.TextMessage, []byte(jsonString)); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
