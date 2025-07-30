package modules

import (
	"testing"

	"github.com/TLop503/LogCrunch/structs"
)

func TestParseApacheLog(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected *structs.ApacheLogEntry
		wantErr  bool
	}{
		{
			name: "Valid log line",
			line: `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`,
			expected: &structs.ApacheLogEntry{
				Remote:     "127.0.0.1",
				RemoteLong: "-",
				RemoteUser: "frank",
				Timestamp:  "10/Oct/2000:13:55:36 -0700",
				Request:    "GET /apache_pb.gif HTTP/1.0",
				StatusCode: "200",
				Size:       "2326",
			},
			wantErr: false,
		},
		{
			name:     "Malformed log line",
			line:     `incomplete log line here`,
			expected: nil,
			wantErr:  true,
		},
		{
			name: "Missing fields",
			line: `192.168.0.1 - - [28/Jul/2025:10:23:45 -0700] "POST /login HTTP/1.1" 401 -`,
			expected: &structs.ApacheLogEntry{
				Remote:     "192.168.0.1",
				RemoteLong: "-",
				RemoteUser: "-",
				Timestamp:  "28/Jul/2025:10:23:45 -0700",
				Request:    "POST /login HTTP/1.1",
				StatusCode: "401",
				Size:       "-",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseApacheLog(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseApacheLog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				if *got != *tt.expected {
					t.Errorf("ParseApacheLog() = %+v, want %+v", got, tt.expected)
				}
			}
		})
	}
}
