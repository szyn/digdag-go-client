package digdag

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetLogFiles(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/files.json")
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, false)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	files, err := client.GetLogFiles("64")
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
	if len(files) != 2 {
		t.Fatalf("result should be one: %d", len(files))
	}
}

func TestGetLogFileResult(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/+test+test@5a54eea130ef7740.73100@test.local.log.gz")
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, false)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	logText, err := client.GetLogText("64", "+test+test@5a54eea130ef7740.73100@test.local.log.gz")
	if !strings.Contains(logText, "test") {
		t.Error("log should be contains test but: ", logText)
	}
}
