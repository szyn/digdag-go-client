package digdag

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClient_GetTasks(t *testing.T) {
	type args struct {
		attemptID string
	}
	tests := []struct {
		name    string
		res     string
		args    args
		want    []*Task
		wantErr bool
	}{
		// Test cases
		{
			args: args{attemptID: "64"},
			res: `
			{
				"tasks": [
					{
						"id": "236",
						"fullName": "+test",
						"parentId": null,
						"config": {},
						"upstreams": [],
						"state": "success",
						"exportParams": {},
						"storeParams": {},
						"stateParams": {},
						"updatedAt": "2018-01-09T16:32:34Z",
						"retryAt": null,
						"startedAt": null,
						"isGroup": true
					},
					{
						"id": "237",
						"fullName": "+test+test1",
						"parentId": "236",
						"config": {
							"echo>": "test"
						},
						"upstreams": [],
						"state": "success",
						"exportParams": {},
						"storeParams": {},
						"stateParams": {},
						"updatedAt": "2018-01-09T16:32:34Z",
						"retryAt": null,
						"startedAt": "2018-01-09T16:32:33Z",
						"isGroup": false
					}
				]
			}			
			`,
			want: []*Task{
				{
					ID:           "236",
					FullName:     "+test",
					ParentID:     nil,
					Config:       map[string]interface{}{},
					Upstreams:    []string{},
					State:        "success",
					ExportParams: map[string]interface{}{},
					StoreParams:  map[string]string{},
					StateParams:  map[string]string{},
					UpdatedAt:    "2018-01-09T16:32:34Z",
					RetryAt:      nil,
					StartedAt:    nil,
					IsGroup:      true,
				},
				{
					ID:           "237",
					FullName:     "+test+test1",
					ParentID:     "236",
					Config:       map[string]interface{}{"echo>": "test"},
					Upstreams:    []string{},
					State:        "success",
					ExportParams: map[string]interface{}{},
					StoreParams:  map[string]string{},
					StateParams:  map[string]string{},
					UpdatedAt:    "2018-01-09T16:32:34Z",
					RetryAt:      nil,
					StartedAt:    "2018-01-09T16:32:33Z",
					IsGroup:      false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantURLPath := fmt.Sprintf("/api/attempts/%s/tasks", tt.args.attemptID)
				if r.URL.Path != wantURLPath {
					t.Errorf("URL Path = %v, want : %v", r.URL.Path, wantURLPath)
				}
				fmt.Fprintln(w, tt.res)
			}))
			defer ts.Close()
			c := newTestClient(ts.URL)
			got, err := c.GetTasks(tt.args.attemptID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetTaskResult(t *testing.T) {
	type args struct {
		attemptIDs []string
		taskName   string
	}
	tests := []struct {
		name    string
		res     string
		args    args
		want    *Task
		wantErr bool
	}{
		// Test cases
		{
			name: "get a task status (success)",
			args: args{attemptIDs: []string{"27"}, taskName: "+test+test1"},
			res: `
			{
				"tasks": [
					{
						"id": "236",
						"fullName": "+test",
						"parentId": null,
						"config": {},
						"upstreams": [],
						"state": "success",
						"exportParams": {},
						"storeParams": {},
						"stateParams": {},
						"updatedAt": "2018-01-09T16:32:34Z",
						"retryAt": null,
						"startedAt": null,
						"isGroup": true
					},
					{
						"id": "237",
						"fullName": "+test+test1",
						"parentId": "236",
						"config": {
							"echo>": "test"
						},
						"upstreams": [],
						"state": "success",
						"exportParams": {},
						"storeParams": {},
						"stateParams": {},
						"updatedAt": "2018-01-09T16:32:34Z",
						"retryAt": null,
						"startedAt": "2018-01-09T16:32:33Z",
						"isGroup": false
					}
				]
			}			
			`,
			want: &Task{
				ID:           "237",
				FullName:     "+test+test1",
				ParentID:     "236",
				Config:       map[string]interface{}{"echo>": "test"},
				Upstreams:    []string{},
				State:        "success",
				ExportParams: map[string]interface{}{},
				StoreParams:  map[string]string{},
				StateParams:  map[string]string{},
				UpdatedAt:    "2018-01-09T16:32:34Z",
				RetryAt:      nil,
				StartedAt:    "2018-01-09T16:32:33Z",
				IsGroup:      false,
			},
		},
		{
			name: "invalid task name (doesn't have prefix `+`)",
			args: args{attemptIDs: []string{"27"}, taskName: "test+test1"},
			res: `
			{
				"tasks": [
					{
						"id": "236",
						"fullName": "+test",
						"parentId": null,
						"config": {},
						"upstreams": [],
						"state": "success",
						"exportParams": {},
						"storeParams": {},
						"stateParams": {},
						"updatedAt": "2018-01-09T16:32:34Z",
						"retryAt": null,
						"startedAt": null,
						"isGroup": true
					},
					{
						"id": "237",
						"fullName": "+test+test1",
						"parentId": "236",
						"config": {
							"echo>": "test"
						},
						"upstreams": [],
						"state": "success",
						"exportParams": {},
						"storeParams": {},
						"stateParams": {},
						"updatedAt": "2018-01-09T16:32:34Z",
						"retryAt": null,
						"startedAt": "2018-01-09T16:32:33Z",
						"isGroup": false
					}
				]
			}			
			`,
			wantErr: true,
		},
		{
			name: "result not found",
			args: args{attemptIDs: []string{"27"}, taskName: "+test+test1000"},
			res: `
			{
				"tasks": [
					{
						"id": "236",
						"fullName": "+test",
						"parentId": null,
						"config": {},
						"upstreams": [],
						"state": "success",
						"exportParams": {},
						"storeParams": {},
						"stateParams": {},
						"updatedAt": "2018-01-09T16:32:34Z",
						"retryAt": null,
						"startedAt": null,
						"isGroup": true
					},
					{
						"id": "237",
						"fullName": "+test+test1",
						"parentId": "236",
						"config": {
							"echo>": "test"
						},
						"upstreams": [],
						"state": "success",
						"exportParams": {},
						"storeParams": {},
						"stateParams": {},
						"updatedAt": "2018-01-09T16:32:34Z",
						"retryAt": null,
						"startedAt": "2018-01-09T16:32:33Z",
						"isGroup": false
					}
				]
			}			
			`,
			wantErr: true,
		},
		{
			name: "task state is `failed`",
			args: args{attemptIDs: []string{"77"}, taskName: "+test+test1"},
			res: `
			{
				"tasks": [
					{
						"id": "276",
						"fullName": "+test",
						"parentId": null,
						"config": {},
						"upstreams": [],
						"state": "group_error",
						"exportParams": {},
						"storeParams": {},
						"stateParams": {},
						"updatedAt": "2018-01-14T13:16:45Z",
						"retryAt": null,
						"startedAt": null,
						"isGroup": true
					},
					{
						"id": "277",
						"fullName": "+test+test1",
						"parentId": "276",
						"config": {
							"sh>": "exit 1"
						},
						"upstreams": [],
						"state": "error",
						"exportParams": {},
						"storeParams": {},
						"stateParams": {},
						"updatedAt": "2018-01-14T13:16:44Z",
						"retryAt": null,
						"startedAt": "2018-01-14T13:16:44Z",
						"isGroup": false
					},
					{
						"id": "278",
						"fullName": "+test^failure-alert",
						"parentId": "276",
						"config": {
							"_type": "notify",
							"_command": "Workflow session attempt failed"
						},
						"upstreams": [],
						"state": "success",
						"exportParams": {},
						"storeParams": {},
						"stateParams": {},
						"updatedAt": "2018-01-14T13:16:45Z",
						"retryAt": null,
						"startedAt": "2018-01-14T13:16:44Z",
						"isGroup": false
					}
				]
			}			
			`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for _, v := range tt.args.attemptIDs {
					wantURLPath := fmt.Sprintf("/api/attempts/%s/tasks", v)
					if r.URL.Path != wantURLPath {
						t.Errorf("URL Path = %v, want : %v", r.URL.Path, wantURLPath)
					}
				}
				fmt.Fprintln(w, tt.res)
			}))
			defer ts.Close()
			c := newTestClient(ts.URL)
			got, err := c.GetTaskResult(tt.args.attemptIDs, tt.args.taskName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetTaskResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetTaskResult() = %v, want %v", got, tt.want)
			}
		})
	}
}
