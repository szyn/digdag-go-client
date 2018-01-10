package digdag

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWorkflows(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/workflows" {
			t.Errorf("request URL should be /api/workflows but : %v", req.URL.Path)
		}
		http.ServeFile(res, req, "testdata/workflows.json")
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, false)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	workflows, err := client.GetWorkflows()
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
	if len(workflows) != 1 {
		t.Fatalf("result should be one: %d", len(workflows))
	}
}

func TestGetWorkflow(t *testing.T) {
	testProjID := "1"
	expectWfID := "18"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/projects/"+testProjID+"/workflows" {
			t.Errorf("request URL should be /api/projects/%v/workflows but : %v", testProjID, req.URL.Path)
		}
		http.ServeFile(res, req, "testdata/workflows.json")
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, false)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	workflow, err := client.GetWorkflow("1", "test")
	if err != nil {
		t.Error("err shoud be nil but: ", err)
	}
	if workflow.ID != expectWfID {
		t.Errorf("got %v, want %v", workflow.ID, expectWfID)
	}
}
