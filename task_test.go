package digdag

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTasks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			t.Errorf("Expected GET request, got '%s'", req.Method)
		}
		if req.URL.Path != "/api/attempts/27/tasks" {
			t.Error("request URL should be /api/attempts/27/tasks but :", req.URL.Path)
		}

		respJSONFile, err := ioutil.ReadFile(`testdata/tasks.json`)
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

	result, err := client.GetTaskResult([]string{"27"}, "+test+setup")
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
	if result.State != "success" {
		t.Fatalf("want %v but %v", "success", result.State)
	}

	result, err = client.GetTaskResult([]string{"27"}, "+test+failed")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
}
