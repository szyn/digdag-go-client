package digdag

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type sessionsWrapper struct {
	Sessions []Session `json:"sessions"`
}

// Session is the struct for digdag session
type Session struct {
	ID      string `json:"id"`
	Project struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"project"`
	Workflow struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"workflow"`
	SessionUUID string    `json:"sessionUuid"`
	SessionTime time.Time `json:"sessionTime"`
	LastAttempt struct {
		ID               string      `json:"id"`
		RetryAttemptName interface{} `json:"retryAttemptName"`
		Done             bool        `json:"done"`
		Success          bool        `json:"success"`
		CancelRequested  bool        `json:"cancelRequested"`
		Params           struct {
		} `json:"params"`
		CreatedAt  time.Time `json:"createdAt"`
		FinishedAt time.Time `json:"finishedAt"`
	}
}

// GetProjectWorkflowSessions to get sessions by projectID and workflow
func (c *Client) GetProjectWorkflowSessions(projectID, workflowName string) ([]Session, error) {
	spath := fmt.Sprintf("/api/projects/%s/sessions", projectID)

	params := url.Values{}
	params.Set("workflow", workflowName)

	var sw *sessionsWrapper
	err := c.doReq(http.MethodGet, spath, params, &sw)
	if err != nil {
		return nil, err
	}

	// if any sessions not found
	if len(sw.Sessions) == 0 {
		return nil, errors.New("any sessions not found")
	}

	return sw.Sessions, err
}
