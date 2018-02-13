package digdag

import (
	"fmt"
	"net/http"
	"strings"
)

type tasksWrapper struct {
	Tasks []*Task `json:"tasks"`
}

// Task is struct for attempts task result
type Task struct {
	ID           string                 `json:"id"`
	FullName     string                 `json:"fullName"`
	ParentID     interface{}            `json:"parentId"`
	Config       map[string]interface{} `json:"config"`
	Upstreams    []string               `json:"upstreams"`
	State        string                 `json:"state"`
	ExportParams map[string]interface{} `json:"exportParams"`
	StoreParams  map[string]interface{} `json:"storeParams"`
	StateParams  map[string]interface{} `json:"stateParams"`
	UpdatedAt    string                 `json:"updatedAt"`
	RetryAt      interface{}            `json:"retryAt"`
	StartedAt    interface{}            `json:"startedAt"`
	IsGroup      bool                   `json:"isGroup"`
}

// GetTasks to get tasks list
func (c *Client) GetTasks(attemptID string) ([]*Task, error) {
	spath := fmt.Sprintf("/api/attempts/%s/tasks", attemptID)

	var tw *tasksWrapper

	resp, err := c.NewRequest(http.MethodGet, spath, nil)
	if err != nil {
		return nil, err
	}

	if err := decodeBody(resp, &tw); err != nil {
		return nil, err
	}

	return tw.Tasks, err
}

// GetTaskResult to return task result
func (c *Client) GetTaskResult(attemptIDs []string, taskName string) (*Task, error) {
	// Check the taskName has prefix `+`
	if !strings.HasPrefix(taskName, "+") {
		return nil, fmt.Errorf("task `%s` is invalid task name", taskName)
	}

	for _, attemptID := range attemptIDs {
		tasks, err := c.GetTasks(attemptID)
		if err != nil {
			return nil, err
		}

		for k := range tasks {
			if tasks[k].FullName == taskName {
				state := tasks[k].State
				if state == "success" {
					return tasks[k], nil
				}

				return nil, fmt.Errorf("task `%s` state is %s", taskName, state)
			}
		}
	}

	return nil, fmt.Errorf("task `%s` result not found", taskName)
}
