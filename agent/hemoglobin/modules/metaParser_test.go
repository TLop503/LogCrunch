package modules

import (
	"testing"

	"github.com/TLop503/LogCrunch/structs"
)

// -- Apache Tests --

func TestMetaParseApache_Valid(t *testing.T) {
	tests := []struct {
		name     string
		logLine  string
		expected structs.ApacheLogEntry
	}{
		{
			name:    "Basic Apache log",
			logLine: `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`,
			expected: structs.ApacheLogEntry{
				Remote:     "127.0.0.1",
				RemoteLong: "-",
				RemoteUser: "frank",
				Timestamp:  "10/Oct/2000:13:55:36 -0700",
				Request:    "GET /apache_pb.gif HTTP/1.0",
				StatusCode: "200",
				Size:       "2326",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var entry structs.ApacheLogEntry
			err := MetaParse(tt.logLine, ApacheRegex, &entry)
			if err != nil {
				t.Fatalf("MetaParse failed: %v", err)
			}
			if entry != tt.expected {
				t.Errorf("Parsed entry doesn't match expected.\nGot: %#v\nWant: %#v", entry, tt.expected)
			}
		})
	}
}

func TestMetaParseApache_Malformed(t *testing.T) {
	logLine := `this is not a valid apache log line`

	var entry structs.ApacheLogEntry
	err := MetaParse(logLine, ApacheRegex, &entry)
	if err == nil {
		t.Error("Expected error for malformed input, got none")
	}
}

// -- Syslog Tests --

func TestMetaParseSyslog_Valid(t *testing.T) {
	tests := []struct {
		name     string
		logLine  string
		expected structs.SyslogEntry
	}{
		{
			name:    "Basic syslog line",
			logLine: `Jul 30 14:17:01 blackwall CRON[620010]: pam_unix(cron:session): session opened for user root(uid=0) by (uid=0)`,
			expected: structs.SyslogEntry{
				Timestamp: "Jul 30 14:17:01",
				Host:      "blackwall",
				Process:   "CRON",
				PID:       "620010",
				Message:   "pam_unix(cron:session): session opened for user root(uid=0) by (uid=0)",
			},
		},
		{
			name:    "Different process name",
			logLine: `Jul 30 14:20:02 myhost sshd[12345]: Accepted password for user from 10.0.0.1 port 22 ssh2`,
			expected: structs.SyslogEntry{
				Timestamp: "Jul 30 14:20:02",
				Host:      "myhost",
				Process:   "sshd",
				PID:       "12345",
				Message:   "Accepted password for user from 10.0.0.1 port 22 ssh2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var entry structs.SyslogEntry
			err := MetaParse(tt.logLine, SyslogRegex, &entry)
			if err != nil {
				t.Fatalf("MetaParse failed: %v", err)
			}
			if entry != tt.expected {
				t.Errorf("Parsed entry doesn't match expected.\nGot: %#v\nWant: %#v", entry, tt.expected)
			}
		})
	}
}

func TestMetaParseSyslog_Malformed(t *testing.T) {
	logLine := `INVALID SYSLOG ENTRY`

	var entry structs.SyslogEntry
	err := MetaParse(logLine, SyslogRegex, &entry)
	if err == nil {
		t.Error("Expected error for malformed syslog input, got none")
	}
}
