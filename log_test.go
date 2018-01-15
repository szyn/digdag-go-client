package digdag

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClient_GetLogFiles(t *testing.T) {
	type args struct {
		attemptID string
	}
	tests := []struct {
		name    string
		res     string
		args    args
		want    []*LogFile
		wantErr bool
	}{
		// Test cases
		{
			args: args{attemptID: "111"},
			res: `
			{
				"files": [
					{
						"fileName": "+test+test2@5a54eea2007a1200.73100@test.local.log.gz",
						"fileSize": 125,
						"taskName": "+test+test2",
						"fileTime": "2018-01-09T16:32:34Z",
						"agentId": "73100@test.local",
						"direct": null
					},
					{
						"fileName": "+test+test@5a54eea130ef7740.73100@test.local.log.gz",
						"fileSize": 124,
						"taskName": "+test+test",
						"fileTime": "2018-01-09T16:32:33Z",
						"agentId": "73100@test.local",
						"direct": null
					}
				]
			}
			`,
			want: []*LogFile{
				{
					FileName: "+test+test2@5a54eea2007a1200.73100@test.local.log.gz",
					FileSize: 125,
					TaskName: "+test+test2",
					FileTime: "2018-01-09T16:32:34Z",
					AgentID:  "73100@test.local",
					Direct:   nil,
				},
				{
					FileName: "+test+test@5a54eea130ef7740.73100@test.local.log.gz",
					FileSize: 124,
					TaskName: "+test+test",
					FileTime: "2018-01-09T16:32:33Z",
					AgentID:  "73100@test.local",
					Direct:   nil,
				},
			},
		},
		{
			name:    "test files is empty",
			args:    args{attemptID: "11111"},
			res:     `{"files": []}`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantURLPath := fmt.Sprintf("/api/logs/%s/files", tt.args.attemptID)
				if r.URL.Path != wantURLPath {
					t.Errorf("URL Path = %v, want : %v", r.URL.Path, wantURLPath)
				}
				fmt.Fprintln(w, tt.res)
			}))
			defer ts.Close()
			c := newTestClient(ts.URL)
			got, err := c.GetLogFiles(tt.args.attemptID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetLogFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetLogFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetLogFileResult(t *testing.T) {
	type args struct {
		attemptID string
		taskName  string
	}
	tests := []struct {
		name    string
		res     string
		args    args
		want    *LogFile
		wantErr bool
	}{
		// Test cases
		{
			args: args{attemptID: "11", taskName: "+test+test"},
			res: `
			{
				"files": [
					{
						"fileName": "+test+test2@5a54eea2007a1200.73100@test.local.log.gz",
						"fileSize": 125,
						"taskName": "+test+test2",
						"fileTime": "2018-01-09T16:32:34Z",
						"agentId": "73100@test.local",
						"direct": null
					},
					{
						"fileName": "+test+test@5a54eea130ef7740.73100@test.local.log.gz",
						"fileSize": 124,
						"taskName": "+test+test",
						"fileTime": "2018-01-09T16:32:33Z",
						"agentId": "73100@test.local",
						"direct": null
					}
				]
			}
			`,
			want: &LogFile{
				FileName: "+test+test@5a54eea130ef7740.73100@test.local.log.gz",
				FileSize: 124,
				TaskName: "+test+test",
				FileTime: "2018-01-09T16:32:33Z",
				AgentID:  "73100@test.local",
				Direct:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantURLPath := fmt.Sprintf("/api/logs/%s/files", tt.args.attemptID)
				if r.URL.Path != wantURLPath {
					t.Errorf("URL Path = %v, want : %v", r.URL.Path, wantURLPath)
				}
				fmt.Fprintln(w, tt.res)
			}))
			defer ts.Close()
			c := newTestClient(ts.URL)
			got, err := c.GetLogFileResult(tt.args.attemptID, tt.args.taskName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetLogFileResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetLogFileResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetLogText(t *testing.T) {
	type args struct {
		attemptID string
		fileName  string
	}
	tests := []struct {
		name    string
		resFile string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				attemptID: "11",
				fileName:  "+test+test@5a54eea130ef7740.73100@test.local.log.gz",
			},
			resFile: "testdata/+test+test@5a54eea130ef7740.73100@test.local.log.gz",
			want:    "2018-01-10 01:32:34.003 +0900 [INFO] (0315@[0:test]+test+test) io.digdag.core.agent.OperatorManager: echo>: test\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantURLPath := fmt.Sprintf("/api/logs/%s/files/%s", tt.args.attemptID, tt.args.fileName)
				if r.URL.Path != wantURLPath {
					t.Errorf("URL Path = %v, want : %v", r.URL.Path, wantURLPath)
				}
				http.ServeFile(w, r, tt.resFile)
			}))
			defer ts.Close()
			c := newTestClient(ts.URL)
			got, err := c.GetLogText(tt.args.attemptID, tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetLogText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.GetLogText() = %v, want %v", got, tt.want)
			}
		})
	}
}
