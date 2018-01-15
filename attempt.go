package digdag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type attemptsWrapper struct {
	Attempts []*Attempt `json:"attempts"`
}

// Attempt is the struct for digdag attempt
type Attempt struct {
	ID      string `json:"id"`
	Index   int    `json:"index"`
	Project struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"project"`
	Workflow struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"workflow"`
	SessionID        string            `json:"sessionId"`
	SessionUUID      string            `json:"sessionUuid"`
	SessionTime      string            `json:"sessionTime"` // FIXME fix type from string to time.Time
	RetryAttemptName interface{}       `json:"retryAttemptName"`
	Done             bool              `json:"done"`
	Success          bool              `json:"success"`
	CancelRequested  bool              `json:"cancelRequested"`
	Params           map[string]string `json:"params"`
	CreatedAt        string            `json:"createdAt"`
	FinishedAt       string            `json:"finishedAt"`
}

// CreateAttempt is struct for create a new attempt
type CreateAttempt struct {
	Attempt
	WorkflowID       string            `json:"workflowId"`
	SessionTime      string            `json:"sessionTime"`
	RetryAttemptName string            `json:"retryAttemptName,omitempty"`
	Params           map[string]string `json:"params"`
}

// NewCreateAttempt to create a new CreateAttempt struct
func NewCreateAttempt(workflowID, sessionTime, retryAttemptName string) *CreateAttempt {
	return &CreateAttempt{
		WorkflowID:       workflowID,
		SessionTime:      sessionTime,
		RetryAttemptName: retryAttemptName,
		Params:           map[string]string{},
	}
}

// GetAttempts get attempts response
func (c *Client) GetAttempts(attempt *Attempt, includeRetried bool) ([]*Attempt, error) {
	spath := "/api/attempts"

	if attempt == nil {
		attempt = new(Attempt)
	}

	project := attempt.Project.Name
	workflow := attempt.Workflow.Name

	var aw *attemptsWrapper
	ro := &RequestOpts{
		Params: map[string]string{
			"project":         project,
			"workflow":        workflow,
			"include_retried": strconv.FormatBool(includeRetried),
		},
	}

	resp, err := c.NewRequest(http.MethodGet, spath, ro)
	if err != nil {
		return nil, err
	}

	if err := decodeBody(resp, &aw); err != nil {
		return nil, err
	}

	// If any attempts not found
	if len(aw.Attempts) == 0 {
		return nil, fmt.Errorf("attempts does not exist. project=%s workflow=%s", project, workflow)
	}

	return aw.Attempts, nil
}

// GetAttemptIDs to get attemptID from sessionTime
func (c *Client) GetAttemptIDs(projectName, workflowName, targetSession string) (attemptIDs []string, err error) {
	params := new(Attempt)
	params.Project.Name = projectName
	params.Workflow.Name = workflowName

	attempts, err := c.GetAttempts(params, true)
	if err != nil {
		return nil, err
	}

	for k := range attempts {
		sessionTime := attempts[k].SessionTime

		if sessionTime == targetSession {
			attemptIDs = append(attemptIDs, attempts[k].ID)
		}
	}

	// If any attemptID not found
	if len(attemptIDs) == 0 {
		return []string{}, fmt.Errorf("attempts does not exist. project=%s workflow=%s sessionTime=%s", projectName, workflowName, targetSession)
	}

	return attemptIDs, nil
}

// CreateNewAttempt to create a new attempt
func (c *Client) CreateNewAttempt(workflowID, sessionTime string, params []string, retry bool) (attempt *Attempt, done bool, err error) {
	spath := "/api/attempts"

	ca := NewCreateAttempt(workflowID, sessionTime, "")

	// Set params
	for _, v := range params {
		if strings.Contains(v, "=") {
			result := strings.Split(v, "=")
			key := result[0]
			val := result[1]
			ca.Params[key] = val
		}
	}

	// Retry workflow
	if retry == true {
		// TODO: add retry attempt name (optional)
		generatedUUID := uuid.NewV4()
		textID, err := generatedUUID.MarshalText()
		if err != nil {
			return nil, false, err
		}
		ca.RetryAttemptName = string(textID)
	}

	// Create new attempt
	body, err := json.Marshal(ca)
	if err != nil {
		return nil, false, err
	}

	ro := &RequestOpts{
		Body: bytes.NewBuffer(body),
	}

	resp, err := c.NewRequest(http.MethodPut, spath, ro)
	if err != nil {
		// if already session exist
		if resp.StatusCode == http.StatusConflict {
			return nil, true, nil
		}

		return nil, false, err
	}

	if err := decodeBody(resp, &attempt); err != nil {
		return nil, false, err
	}

	return attempt, done, err
}
