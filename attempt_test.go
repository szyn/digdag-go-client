package digdag

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClient_GetAttempts(t *testing.T) {
	type args struct {
		attempt        *Attempt
		includeRetried bool
	}
	tests := []struct {
		name    string
		args    args
		res     string
		want    []*Attempt
		wantErr bool
	}{
		// Test cases
		{
			args: args{
				attempt: &Attempt{
					Project: struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					}{
						ID: "1",
					},
					Workflow: struct {
						Name string `json:"name"`
						ID   string `json:"id"`
					}{
						Name: "test",
					},
				},
				includeRetried: true,
			},
			res: `
			{
				"attempts": [
					  {
						"id": "27",
						"index": 1,
						"project": {
							"id": "1",
							"name": "test"
						},
						"workflow": {
							"name": "test",
							"id": "2"
						},
						"sessionId": "9",
						"sessionUuid": "15624750-5c1f-45d2-b668-c4f86e757484",
						"sessionTime": "2017-06-24T00:00:00+00:00",
						"retryAttemptName": null,
						"done": true,
						"success": false,
						"cancelRequested": false,
						"params": {},
						"createdAt": "2017-06-24T06:45:26Z",
						"finishedAt": "2017-06-24T06:45:31Z"
					}
				]
			}
			`,
			want: []*Attempt{
				{
					ID:    "27",
					Index: 1,
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
					SessionID:        "9",
					SessionUUID:      "15624750-5c1f-45d2-b668-c4f86e757484",
					SessionTime:      "2017-06-24T00:00:00+00:00",
					RetryAttemptName: nil,
					Done:             true,
					Success:          false,
					CancelRequested:  false,
					Params:           map[string]string{},
					CreatedAt:        "2017-06-24T06:45:26Z",
					FinishedAt:       "2017-06-24T06:45:31Z",
				},
			},
		},
		{
			name: "test args attempts nil",
			args: args{
				includeRetried: true,
			},
			res: `
			{
				"attempts": [
					  {
						"id": "27",
						"index": 1,
						"project": {
							"id": "1",
							"name": "test"
						},
						"workflow": {
							"name": "test",
							"id": "2"
						},
						"sessionId": "9",
						"sessionUuid": "15624750-5c1f-45d2-b668-c4f86e757484",
						"sessionTime": "2017-06-24T00:00:00+00:00",
						"retryAttemptName": null,
						"done": true,
						"success": false,
						"cancelRequested": false,
						"params": {},
						"createdAt": "2017-06-24T06:45:26Z",
						"finishedAt": "2017-06-24T06:45:31Z"
					}
				]
			}
			`,
			want: []*Attempt{
				{
					ID:    "27",
					Index: 1,
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
					SessionID:        "9",
					SessionUUID:      "15624750-5c1f-45d2-b668-c4f86e757484",
					SessionTime:      "2017-06-24T00:00:00+00:00",
					RetryAttemptName: nil,
					Done:             true,
					Success:          false,
					CancelRequested:  false,
					Params:           map[string]string{},
					CreatedAt:        "2017-06-24T06:45:26Z",
					FinishedAt:       "2017-06-24T06:45:31Z",
				},
			},
		},
		{
			args: args{
				attempt: &Attempt{
					Project: struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					}{
						ID: "11111",
					},
					Workflow: struct {
						Name string `json:"name"`
						ID   string `json:"id"`
					}{
						Name: "hoge",
					},
				},
				includeRetried: true,
			},
			res:     `{"attempts": []}`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantURLPath := "/api/attempts"
				if r.URL.Path != wantURLPath {
					t.Errorf("URL Path = %v, want : %v", r.URL.Path, wantURLPath)
				}
				fmt.Fprintln(w, tt.res)
			}))
			defer ts.Close()
			c := newTestClient(ts.URL)
			got, err := c.GetAttempts(tt.args.attempt, tt.args.includeRetried)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetAttempts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetAttempts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetAttemptIDs(t *testing.T) {
	type args struct {
		projectName   string
		workflowName  string
		targetSession string
	}
	tests := []struct {
		name           string
		res            string
		args           args
		wantAttemptIDs []string
		wantErr        bool
	}{
		// Test cases
		{
			args: args{
				projectName:   "test",
				workflowName:  "test",
				targetSession: "2017-06-24T00:00:00+00:00",
			},
			res: `
			{
				"attempts": [
					  {
						"id": "27",
						"index": 1,
						"project": {
							"id": "1",
							"name": "test"
						},
						"workflow": {
							"name": "test",
							"id": "2"
						},
						"sessionId": "9",
						"sessionUuid": "15624750-5c1f-45d2-b668-c4f86e757484",
						"sessionTime": "2017-06-24T00:00:00+00:00",
						"retryAttemptName": null,
						"done": true,
						"success": false,
						"cancelRequested": false,
						"params": {},
						"createdAt": "2017-06-24T06:45:26Z",
						"finishedAt": "2017-06-24T06:45:31Z"
					}
				]
			}			
			`,
			wantAttemptIDs: []string{"27"},
		},
		{
			args: args{
				projectName:   "test",
				workflowName:  "test",
				targetSession: "2017-06-25T00:00:00+00:00",
			},
			res: `
			{
				"attempts": [
					  {
						"id": "27",
						"index": 1,
						"project": {
							"id": "1",
							"name": "test"
						},
						"workflow": {
							"name": "test",
							"id": "2"
						},
						"sessionId": "9",
						"sessionUuid": "15624750-5c1f-45d2-b668-c4f86e757484",
						"sessionTime": "2017-06-24T00:00:00+00:00",
						"retryAttemptName": null,
						"done": true,
						"success": false,
						"cancelRequested": false,
						"params": {},
						"createdAt": "2017-06-24T06:45:26Z",
						"finishedAt": "2017-06-24T06:45:31Z"
					}
				]
			}			
			`,
			wantAttemptIDs: []string{},
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantURLPath := "/api/attempts"
				if r.URL.Path != wantURLPath {
					t.Errorf("URL Path = %v, want : %v", r.URL.Path, wantURLPath)
				}
				fmt.Fprintln(w, tt.res)
			}))
			defer ts.Close()
			c := newTestClient(ts.URL)
			gotAttemptIDs, err := c.GetAttemptIDs(tt.args.projectName, tt.args.workflowName, tt.args.targetSession)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetAttemptIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAttemptIDs, tt.wantAttemptIDs) {
				t.Errorf("Client.GetAttemptIDs() = %v, want %v", gotAttemptIDs, tt.wantAttemptIDs)
			}
		})
	}
}

func TestNewCreateAttempt(t *testing.T) {
	type args struct {
		workflowID       string
		sessionTime      string
		retryAttemptName string
	}
	tests := []struct {
		name string
		args args
		want *CreateAttempt
	}{
		// Test cases
		{
			args: args{
				workflowID:       "2",
				sessionTime:      "2017-06-24T00:00:00+00:00",
				retryAttemptName: "",
			},
			want: &CreateAttempt{
				WorkflowID:       "2",
				SessionTime:      "2017-06-24T00:00:00+00:00",
				RetryAttemptName: "",
				Params:           map[string]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCreateAttempt(tt.args.workflowID, tt.args.sessionTime, tt.args.retryAttemptName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCreateAttempt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_CreateNewAttempt(t *testing.T) {
	type args struct {
		workflowID  string
		sessionTime string
		params      []string
		retry       bool
	}
	tests := []struct {
		name        string
		res         string
		args        args
		wantAttempt *Attempt
		wantDone    bool
		wantErr     bool
	}{
		// Test cases
		{
			name: "test start a attempt",
			args: args{
				workflowID:  "2",
				sessionTime: "2017-06-24T00:00:00+00:00",
				params:      []string{"key=value"},
				retry:       false,
			},
			res: `
			{
				"id": "27",
				"index": 1,
				"project": {
					"id": "1",
					"name": "test"
				},
				"workflow": {
					"name": "test",
					"id": "2"
				},
				"sessionId": "9",
				"sessionUuid": "15624750-5c1f-45d2-b668-c4f86e757484",
				"done": false,
				"success": false,
				"cancelRequested": false,
				"createdAt": "2017-06-24T06:45:26Z",
				"finishedAt": "",
				"workflowId": "2",
				"sessionTime": "2017-06-24T00:00:00+00:00",
				"params": {
					"key": "value"
				}
			}	
			`,
			wantAttempt: &Attempt{
				ID:    "27",
				Index: 1,
				Project: struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				}{
					"1", "test",
				},
				Workflow: struct {
					Name string `json:"name"`
					ID   string `json:"id"`
				}{
					"test", "2",
				},
				SessionID:       "9",
				SessionUUID:     "15624750-5c1f-45d2-b668-c4f86e757484",
				Done:            false,
				Success:         false,
				CancelRequested: false,
				CreatedAt:       "2017-06-24T06:45:26Z",
				Params:          map[string]string{"key": "value"},
				SessionTime:     "2017-06-24T00:00:00+00:00",
			},
		},
		{
			name: "test if already a session has done",
			args: args{
				workflowID:  "2",
				sessionTime: "2017-06-24T00:00:00+00:00",
				params:      []string{},
				retry:       true,
			},
			wantDone: true,
		},
		{
			name: "test retry a attempt",
			args: args{
				workflowID:  "2",
				sessionTime: "2017-06-24T00:00:00+00:00",
				params:      []string{},
				retry:       true,
			},
			res: `
			{
				"id": "27",
				"index": 1,
				"project": {
					"id": "1",
					"name": "test"
				},
				"workflow": {
					"name": "test",
					"id": "2"
				},
				"sessionId": "9",
				"sessionUuid": "15624750-5c1f-45d2-b668-c4f86e757484",
				"done": false,
				"success": false,
				"cancelRequested": false,
				"createdAt": "2017-06-24T06:45:26Z",
				"finishedAt": "",
				"workflowId": "2",
				"sessionTime": "2017-06-24T00:00:00+00:00",
				"params": {}
			}
			`,
			wantAttempt: &Attempt{
				ID:    "27",
				Index: 1,
				Project: struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				}{
					"1", "test",
				},
				Workflow: struct {
					Name string `json:"name"`
					ID   string `json:"id"`
				}{
					"test", "2",
				},
				SessionID:       "9",
				SessionUUID:     "15624750-5c1f-45d2-b668-c4f86e757484",
				Done:            false,
				Success:         false,
				CancelRequested: false,
				CreatedAt:       "2017-06-24T06:45:26Z",
				Params:          map[string]string{},
				SessionTime:     "2017-06-24T00:00:00+00:00",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantURLPath := "/api/attempts"
				if r.URL.Path != wantURLPath {
					t.Errorf("URL Path = %v, want : %v", r.URL.Path, wantURLPath)
				}
				// Case StatusConflict
				if tt.name == "test if already a session has done" {
					w.WriteHeader(http.StatusConflict)
				}
				fmt.Fprintln(w, tt.res)
			}))
			defer ts.Close()
			c := newTestClient(ts.URL)
			gotAttempt, gotDone, err := c.CreateNewAttempt(tt.args.workflowID, tt.args.sessionTime, tt.args.params, tt.args.retry)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateNewAttempt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAttempt, tt.wantAttempt) {
				t.Errorf("Client.CreateNewAttempt() gotAttempt = %v, want %v", gotAttempt, tt.wantAttempt)
			}
			if gotDone != tt.wantDone {
				t.Errorf("Client.CreateNewAttempt() gotDone = %v, want %v", gotDone, tt.wantDone)
			}
		})
	}
}
