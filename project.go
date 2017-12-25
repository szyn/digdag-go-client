package digdag

import (
	"errors"
	"net/http"
	"net/url"
)

// projectsWrapper is struct for received json
type projectsWrapper struct {
	Projects []Project `json:"projects"`
}

// Project is struct for digdag project
type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//
// func (c *Client) GetProjects() ([]Project, error) {
// 	spath := "/api/projects"

// 	var pw *projectsWrapper
// 	err := c.doReq(http.MethodGet, spath, nil, &pw)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return pw.Projects, nil
// }

func (c *Client) GetProjects() ([]Project, error) {
	spath := "/api/projects"

	var pw *projectsWrapper
	resp, err := c.NewRequest(http.MethodGet, spath, nil)
	if err != nil {
		return nil, err
	}
	decodeBody(resp, &pw)

	return pw.Projects, nil
}

// GetProjectIDByName to get project ID by project name
func (c *Client) GetProjectIDByName(name string) (projectID string, err error) {
	spath := "/api/projects"

	params := url.Values{}
	// params.Set("name", c.ProjectName)
	params.Set("name", name)

	var pw *projectsWrapper

	ro := &RequestOpts{
		Params: map[string]string{
			"name": name,
		},
	}

	resp, err := c.NewRequest(http.MethodGet, spath, ro)
	// err = c.doReq(http.MethodGet, spath, params, &pw)
	// err = c.doReq(http.MethodGet, spath, params, &projects)
	if err != nil {
		return "", err
	}

	decodeBody(resp, &pw)

	// if project not found
	if len(pw.Projects) == 0 {
		return "", errors.New("project not found: `" + name + "`")

		// return "", errors.New("project not found: `" + c.ProjectName + "`")
	}

	projectID = pw.Projects[0].ID

	return projectID, nil
}
