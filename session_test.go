package digdag

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestClient_GetSessions(t *testing.T) {
	sessionTime, _ := time.Parse("2006-01-02T15:04:05-07:00", "2017-10-08T04:37:42+00:00")
	createdAt, _ := time.Parse("2006-01-02T15:04:05Z", "2017-10-08T04:37:42Z")
	finishedAt, _ := time.Parse("2006-01-02T15:04:05Z", "2017-10-08T04:37:44Z")

	tests := []struct {
		name    string
		res     string
		want    []*Session
		wantErr bool
	}{
		// Test cases
		{
			res: `
			{
				"sessions": [
					{
						"id": "2",
						"project": {
							"id": "1",
							"name": "test"
						},
						"workflow": {
							"name": "test",
							"id": "2"
						},
						"sessionUuid": "eaf514b8-b40b-4aea-81e4-9f46c0e2d3d5",
						"sessionTime": "2017-10-08T04:37:42+00:00",
						"lastAttempt": {
							"id": "2",
							"retryAttemptName": null,
							"done": true,
							"success": false,
							"cancelRequested": false,
							"params": {},
							"createdAt": "2017-10-08T04:37:42Z",
							"finishedAt": "2017-10-08T04:37:44Z"
						}
					}
				]
			}
			`,
			want: []*Session{
				{
					ID: "2",
					Project: struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					}{
						ID:   "1",
						Name: "test",
					},
					Workflow: struct {
						Name string `json:"name"`
						ID   string `json:"id"`
					}{
						Name: "test",
						ID:   "2",
					},
					SessionUUID: "eaf514b8-b40b-4aea-81e4-9f46c0e2d3d5",
					SessionTime: sessionTime,
					LastAttempt: struct {
						ID               string            `json:"id"`
						RetryAttemptName interface{}       `json:"retryAttemptName"`
						Done             bool              `json:"done"`
						Success          bool              `json:"success"`
						CancelRequested  bool              `json:"cancelRequested"`
						Params           map[string]string `json:"params"`
						CreatedAt        time.Time         `json:"createdAt"`
						FinishedAt       time.Time         `json:"finishedAt"`
					}{
						ID:               "2",
						RetryAttemptName: nil,
						Done:             true,
						Success:          false,
						CancelRequested:  false,
						Params:           map[string]string{},
						CreatedAt:        createdAt,
						FinishedAt:       finishedAt,
					},
				},
			},
		},
		// {
		// 	res:     `{"sessions": []}`,
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantURLPath := "/api/sessions"
				if r.URL.Path != wantURLPath {
					t.Errorf("URL Path = %v, want : %v", r.URL.Path, wantURLPath)
				}
				fmt.Fprintln(w, tt.res)
			}))
			defer ts.Close()
			c := newTestClient(ts.URL)
			got, err := c.GetSessions()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetSessions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetSessions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetProjectWorkflowSessions(t *testing.T) {
	sessionTime, _ := time.Parse("2006-01-02T15:04:05-07:00", "2017-10-08T04:37:42+00:00")
	createdAt, _ := time.Parse("2006-01-02T15:04:05Z", "2017-10-08T04:37:42Z")
	finishedAt, _ := time.Parse("2006-01-02T15:04:05Z", "2017-10-08T04:37:44Z")
	type args struct {
		projectID    string
		workflowName string
	}
	tests := []struct {
		name    string
		args    args
		res     string
		want    []*Session
		wantErr bool
	}{
		// Test cases
		{
			args: args{projectID: "1", workflowName: "test"},
			res: `
			{
				"sessions": [
					{
						"id": "2",
						"project": {
							"id": "1",
							"name": "test"
						},
						"workflow": {
							"name": "test",
							"id": "2"
						},
						"sessionUuid": "eaf514b8-b40b-4aea-81e4-9f46c0e2d3d5",
						"sessionTime": "2017-10-08T04:37:42+00:00",
						"lastAttempt": {
							"id": "2",
							"retryAttemptName": null,
							"done": true,
							"success": false,
							"cancelRequested": false,
							"params": {},
							"createdAt": "2017-10-08T04:37:42Z",
							"finishedAt": "2017-10-08T04:37:44Z"
						}
					}
				]
			}
			`,
			want: []*Session{
				{
					ID: "2",
					Project: struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					}{
						ID:   "1",
						Name: "test",
					},
					Workflow: struct {
						Name string `json:"name"`
						ID   string `json:"id"`
					}{
						Name: "test",
						ID:   "2",
					},
					SessionUUID: "eaf514b8-b40b-4aea-81e4-9f46c0e2d3d5",
					SessionTime: sessionTime,
					LastAttempt: struct {
						ID               string            `json:"id"`
						RetryAttemptName interface{}       `json:"retryAttemptName"`
						Done             bool              `json:"done"`
						Success          bool              `json:"success"`
						CancelRequested  bool              `json:"cancelRequested"`
						Params           map[string]string `json:"params"`
						CreatedAt        time.Time         `json:"createdAt"`
						FinishedAt       time.Time         `json:"finishedAt"`
					}{
						ID:               "2",
						RetryAttemptName: nil,
						Done:             true,
						Success:          false,
						CancelRequested:  false,
						Params:           map[string]string{},
						CreatedAt:        createdAt,
						FinishedAt:       finishedAt,
					},
				},
			},
		},
		{
			args:    args{projectID: "33", workflowName: "hoge"},
			res:     `{"sessions":[]}`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantURLPath := fmt.Sprintf("/api/projects/%s/sessions", tt.args.projectID)
				if r.URL.Path != wantURLPath {
					t.Errorf("URL Path = %v, want : %v", r.URL.Path, wantURLPath)
				}
				fmt.Fprintln(w, tt.res)
			}))
			defer ts.Close()
			c := newTestClient(ts.URL)
			got, err := c.GetProjectWorkflowSessions(tt.args.projectID, tt.args.workflowName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetProjectWorkflowSessions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetProjectWorkflowSessions() = %v, want %v", got, tt.want)
			}
		})
	}
}
