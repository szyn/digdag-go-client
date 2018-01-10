package digdag

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

type sessionsWrapper struct {
	Sessions []*Session `json:"sessions"`
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

// GetSessions to get sessions
func (c *Client) GetSessions() ([]*Session, error) {
	spath := "/api/sessions"

	var sw *sessionsWrapper
	resp, err := c.NewRequest(http.MethodGet, spath, nil)
	if err != nil {
		return nil, err
	}

	if err := decodeBody(resp, &sw); err != nil {
		return nil, err
	}

	return sw.Sessions, nil
}

// GetProjectWorkflowSessions to get sessions by projectID and workflow
func (c *Client) GetProjectWorkflowSessions(projectID, workflowName string) ([]*Session, error) {
	spath := fmt.Sprintf("/api/projects/%s/sessions", projectID)

	var sw *sessionsWrapper
	ro := &RequestOpts{
		Params: map[string]string{
			"workflow": workflowName,
		},
	}

	resp, err := c.NewRequest(http.MethodGet, spath, ro)
	if err != nil {
		return nil, err
	}

	if err := decodeBody(resp, &sw); err != nil {
		return nil, err
	}

	// if any sessions not found
	if len(sw.Sessions) == 0 {
		return nil, errors.New("any sessions not found")
	}

	return sw.Sessions, err
}
