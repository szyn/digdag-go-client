package digdag

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClient_GetProjects(t *testing.T) {
	tests := []struct {
		name    string
		res     string
		want    []*Project
		wantErr bool
	}{
		// Test cases
		{
			res: `
			{
				"projects": [
					{
						"id": "1",
						"name": "test",
						"revision": "59b4f254-dc4f-429d-b0d6-18a676bb9e5f",
						"createdAt": "2017-06-13T03:11:12Z",
						"updatedAt": "2017-06-13T03:11:12Z",
						"deletedAt": null,
						"archiveType": "db",
						"archiveMd5": "PeWLToy+/ygCFXQdXlgUbQ=="
					}
				]
			}
			`,
			want: []*Project{
				{
					ID:   "1",
					Name: "test",
				},
			},
		},
		// {
		// 	res:     `{"projects":[]}`,
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantURLPath := "/api/projects"
				if r.URL.Path != wantURLPath {
					t.Errorf("URL Path = %v, want : %v", r.URL.Path, wantURLPath)
				}
				fmt.Fprintln(w, tt.res)
			}))
			defer ts.Close()

			c, err := NewClient(ts.URL, false)
			if err != nil {
				t.Error("err should be nil but: ", err)
			}

			got, err := c.GetProjects()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetProjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetProjects() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetProject(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		res     string
		args    args
		want    *Project
		wantErr bool
	}{
		// Test cases
		{
			args: args{name: "test"},
			res: `
			{
				"projects": [
					{
						"id": "1",
						"name": "test",
						"revision": "59b4f254-dc4f-429d-b0d6-18a676bb9e5f",
						"createdAt": "2017-06-13T03:11:12Z",
						"updatedAt": "2017-06-13T03:11:12Z",
						"deletedAt": null,
						"archiveType": "db",
						"archiveMd5": "PeWLToy+/ygCFXQdXlgUbQ=="
					}
				]
			}
			`,
			want: &Project{
				ID:   "1",
				Name: "test",
			},
		},
		{
			res:     `{"projects":[]}`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantURLPath := "/api/projects"
				if r.URL.Path != wantURLPath {
					t.Errorf("URL Path = %v, want : %v", r.URL.Path, wantURLPath)
				}
				fmt.Fprintln(w, tt.res)
			}))
			defer ts.Close()
			c := newTestClient(ts.URL)
			got, err := c.GetProject(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetProject() = %v, want %v", got, tt.want)
			}
		})
	}
}
