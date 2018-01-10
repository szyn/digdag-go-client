package digdag

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAttempts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			t.Errorf("Expected GET request, got '%s'", req.Method)
		}
		if req.URL.Path != "/api/attempts" {
			t.Error("request URL should be /api/attempts but :", req.URL.Path)
		}

		respJSONFile, err := ioutil.ReadFile(`testdata/attempts.json`)
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

	attempts, err := client.GetAttempts(nil, true)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
	if len(attempts) != 3 {
		t.Fatalf("result should be three: %d", len(attempts))
	}
	if attempts[0].ID != "27" {
		t.Fatalf("want %v but %v", "27", attempts[0].ID)
	}
}

func TestGetAttemptIDs(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			t.Errorf("Expected GET request, got '%s'", req.Method)
		}
		if req.URL.Path != "/api/attempts" {
			t.Error("request URL should be /api/attempts but :", req.URL.Path)
		}

		respJSONFile, err := ioutil.ReadFile(`testdata/attempts.json`)
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

	IDs, err := client.GetAttemptIDs("test", "test", "2017-06-24T00:00:00+00:00")
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
	if len(IDs) != 1 {
		t.Fatalf("result should be one: %d", len(IDs))
	}
	if IDs[0] != "27" {
		t.Fatalf("want %v but %v", "27", IDs[0])
	}
}

func TestCreateNewAttempt(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got '%s'", req.Method)
		}
		if req.URL.Path != "/api/attempts" {
			t.Error("request URL should be /api/attempts but :", req.URL.Path)
		}

		respJSONFile, err := ioutil.ReadFile(`testdata/new_attempt.json`)
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

	result, done, err := client.CreateNewAttempt("2", "2017-06-24", nil, false)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
	if done != false {
		t.Fatalf("want %v but %v", true, result.Done)
	}

	result, done, err = client.CreateNewAttempt("2", "2017-06-24", nil, true)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
	if done != false {
		t.Fatalf("want %v but %v", true, result.Done)
	}

	// Add params
	params := []string{"key1=val1", "key2=val2"}
	result, done, err = client.CreateNewAttempt("2", "2017-06-24", params, true)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
	if done != false {
		t.Fatalf("want %v but %v", true, result.Done)
	}

}
