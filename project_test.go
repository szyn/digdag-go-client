package digdag

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProjects(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			t.Errorf("Expected GET request, got '%s'", req.Method)
		}
		if req.URL.Path != "/api/projects" {
			t.Error("request URL should be /api/projects but :", req.URL.Path)
		}

		respJSONFile, err := ioutil.ReadFile(`testdata/projects.json`)
		if err != nil {
			t.Error("unexpected error: ", err)
		}

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSONFile))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, false)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	// GetProjects
	projects, err := client.GetProjects()
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
	if len(projects) != 1 {
		t.Fatalf("result should be one: %d", len(projects))
	}
	if projects[0].ID != "1" {
		t.Fatalf("want %v but %v", "1", projects[0].ID)
	}

	// GetProject
	project, err := client.GetProject("test")
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
	if project.ID != "1" {
		t.Fatalf("want %v but %v", "1", project.ID)
	}
}
