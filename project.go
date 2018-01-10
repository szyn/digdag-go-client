package digdag

import (
	"fmt"
	"net/http"
)

// projectsWrapper is struct for received json
type projectsWrapper struct {
	Projects []*Project `json:"projects"`
}

// Project is struct for digdag project
type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetProjects to get projects
func (c *Client) GetProjects() ([]*Project, error) {
	spath := "/api/projects"

	var pw *projectsWrapper
	resp, err := c.NewRequest(http.MethodGet, spath, nil)
	if err != nil {
		return nil, err
	}

	if err := decodeBody(resp, &pw); err != nil {
		return nil, err
	}

	return pw.Projects, nil
}

// GetProject to get project by project name
func (c *Client) GetProject(name string) (*Project, error) {
	spath := "/api/projects"

	var pw *projectsWrapper
	ro := &RequestOpts{
		Params: map[string]string{
			"name": name,
		},
	}

	resp, err := c.NewRequest(http.MethodGet, spath, ro)
	if err != nil {
		return nil, err
	}

	if err := decodeBody(resp, &pw); err != nil {
		return nil, err
	}

	// if an empty array (= project not found)
	if len(pw.Projects) == 0 {
		return nil, fmt.Errorf("project `%s` not found", name)
	}

	return pw.Projects[0], nil
}
