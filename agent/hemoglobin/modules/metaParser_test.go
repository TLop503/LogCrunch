package modules

import (
	"github.com/TLop503/LogCrunch/structs"
	"testing"
)

func TestMetaParseApache(t *testing.T) {
	logLine := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`

	var entry structs.ApacheLogEntry

	err := MetaParse(logLine, apacheRegex, &entry)
	if err != nil {
		t.Fatalf("MetaParse failed: %v", err)
	}

	expected := structs.ApacheLogEntry{
		Remote:     "127.0.0.1",
		RemoteLong: "-",
		RemoteUser: "frank",
		Timestamp:  "10/Oct/2000:13:55:36 -0700",
		Request:    "GET /apache_pb.gif HTTP/1.0",
		StatusCode: "200",
		Size:       "2326",
	}

	if entry != expected {
		t.Errorf("Parsed entry doesn't match expected.\nGot: %#v\nWant: %#v", entry, expected)
	}
}

func TestMetaParseSyslog(t *testing.T) {
	logLine := `Jul 30 14:17:01 blackwall CRON[620010]: pam_unix(cron:session): session opened for user root(uid=0) by (uid=0)`

	var entry structs.SyslogEntry

	err := MetaParse(logLine, syslogRegex, &entry)
	if err != nil {
		t.Fatalf("MetaParse failed: %v", err)
	}

	expected := structs.SyslogEntry{
		Timestamp: "Jul 30 14:17:01",
		Host:      "blackwall",
		Process:   "CRON",
		PID:       "620010",
		Message:   "pam_unix(cron:session): session opened for user root(uid=0) by (uid=0)",
	}

	if entry != expected {
		t.Errorf("Parsed entry doesn't match expected.\nGot: %#v\nWant: %#v", entry, expected)
	}
}
