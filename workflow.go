package digdag

import (
	"fmt"
	"net/http"
)

type workflowsWrapper struct {
	Workflows []*Workflow `json:"workflows"`
}

// Workflow is struct for digdag workflow
type Workflow struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Project  `json:"project"`
	Revision string `json:"revision"`
	Timezone string `json:"timezone"`
}

// GetWorkflows to get projects
func (c *Client) GetWorkflows() ([]*Workflow, error) {
	spath := "/api/workflows"

	var ww *workflowsWrapper
	resp, err := c.NewRequest(http.MethodGet, spath, nil)
	if err != nil {
		return nil, err
	}

	if err := decodeBody(resp, &ww); err != nil {
		return nil, err
	}

	return ww.Workflows, nil
}

// GetWorkflow to get workflow by project ID and workflow name
func (c *Client) GetWorkflow(projectID, workflowName string) (*Workflow, error) {
	spath := fmt.Sprintf("/api/projects/%s/workflows", projectID)

	var ww *workflowsWrapper
	ro := &RequestOpts{
		Params: map[string]string{
			"name": workflowName,
		},
	}

	resp, err := c.NewRequest(http.MethodGet, spath, ro)
	if err != nil {
		return nil, err
	}

	if err := decodeBody(resp, &ww); err != nil {
		return nil, err
	}

	// if an empty array (= workflow not found)
	if len(ww.Workflows) == 0 {
		return nil, fmt.Errorf("workflow `%s` not found", workflowName)
	}

	return ww.Workflows[0], nil
}
