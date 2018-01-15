package digdag

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClient_TestGetWorkflows(t *testing.T) {
	tests := []struct {
		name    string
		res     string
		want    []*Workflow
		wantErr bool
	}{
		// Test cases
		{
			res: `
			{
				"workflows": [
					{
						"id": "9",
						"name": "test",
						"project": {
							"id": "3",
							"name": "test"
						},
						"revision": "2c9144e6-4d77-471b-baf6-f7d46f1b5296",
						"timezone": "UTC",
						"config": {
							"+test": {
								"echo>": "test"
							}
						}
					}
				]
			}			
			`,
			want: []*Workflow{
				{
					ID:   "9",
					Name: "test",
					Project: Project{
						ID:   "3",
						Name: "test",
					},
					Revision: "2c9144e6-4d77-471b-baf6-f7d46f1b5296",
					Timezone: "UTC",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantURLPath := "/api/workflows"
				if r.URL.Path != wantURLPath {
					t.Errorf("URL Path = %v, want : %v", r.URL.Path, wantURLPath)
				}
				fmt.Fprintln(w, tt.res)
			}))
			defer ts.Close()
			c := newTestClient(ts.URL)
			got, err := c.GetWorkflows()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetWorkflows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetWorkflows() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_TestGetWorkflow(t *testing.T) {
	type args struct {
		projectID    string
		workflowName string
	}
	tests := []struct {
		name    string
		args    args
		res     string
		want    *Workflow
		wantErr bool
	}{
		// Test cases
		{
			args: args{
				projectID:    "9",
				workflowName: "test",
			},
			res: `
			{
				"workflows": [
					{
						"id": "9",
						"name": "test",
						"project": {
							"id": "3",
							"name": "test"
						},
						"revision": "2c9144e6-4d77-471b-baf6-f7d46f1b5296",
						"timezone": "UTC",
						"config": {
							"+test": {
								"echo>": "test"
							}
						}
					}
				]
			}
			`,
			want: &Workflow{
				ID:   "9",
				Name: "test",
				Project: Project{
					ID:   "3",
					Name: "test",
				},
				Revision: "2c9144e6-4d77-471b-baf6-f7d46f1b5296",
				Timezone: "UTC",
			},
		},
		{
			args: args{
				projectID:    "9",
				workflowName: "hoge",
			},
			res:     `{"workflows":[]}`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantURLPath := fmt.Sprintf("/api/projects/%s/workflows", tt.args.projectID)
				if r.URL.Path != wantURLPath {
					t.Errorf("URL Path = %v, want : %v", r.URL.Path, wantURLPath)
				}
				fmt.Fprintln(w, tt.res)
			}))
			defer ts.Close()
			c := newTestClient(ts.URL)
			got, err := c.GetWorkflow(tt.args.projectID, tt.args.workflowName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetWorkflow() = %v, want %v", got, tt.want)
			}
		})
	}
}
