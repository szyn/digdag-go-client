package digdag

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	parsedURL, _ := url.Parse(defaultBaseURL)
	type args struct {
		urlStr  string
		verbose bool
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{
			name: "Test NewClient",
			args: args{"", false},
			want: &Client{
				BaseURL:       parsedURL,
				HTTPClient:    &http.Client{},
				UserAgent:     defaultUserAgent,
				CustomHeaders: http.Header{},
				Verbose:       false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.urlStr, tt.args.verbose)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
