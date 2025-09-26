package modules

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/TLop503/LogCrunch/structs"
)

// -- Helpers to generate ParserModule for custom tests --

func makeSyslogModule() structs.ParserModule {
	return structs.ParserModule{
		Regex:  regexp.MustCompile(`^(?P<timestamp>\w+\s+\d+\s+\d+:\d+:\d+)\s+(?P<host>\S+)\s+(?P<process>\w+)(?:\[(?P<pid>\d+)\])?:\s+(?P<message>.*)$`),
		Schema: structs.ReflectSchema(structs.SyslogEntry{}),
	}
}

func makeApacheModule() structs.ParserModule {
	return structs.ParserModule{
		Regex:  regexp.MustCompile(`(?P<remote>\S+) (?P<remote_long>\S+) (?P<remote_user>\S+) \[(?P<timestamp>[^\]]+)\] "(?P<request>[^"]*)" (?P<status_code>\d{3}) (?P<size>\S+)`),
		Schema: structs.ReflectSchema(structs.ApacheLogEntry{}),
	}
}

// -- Syslog Tests --

func TestMetaParseSyslog_Registry(t *testing.T) {
	module := structs.MetaParserRegistry["syslog"]

	tests := []struct {
		name     string
		logLine  string
		expected map[string]interface{}
	}{
		{
			name:    "Basic syslog",
			logLine: `Jul 30 14:17:01 blackwall CRON[620010]: pam_unix(cron:session): session opened for user root(uid=0) by (uid=0)`,
			expected: map[string]interface{}{
				"timestamp": "Jul 30 14:17:01",
				"host":      "blackwall",
				"process":   "CRON",
				"pid":       "620010",
				"message":   "pam_unix(cron:session): session opened for user root(uid=0) by (uid=0)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := MetaParse(tt.logLine, module)
			if err != nil {
				t.Fatalf("MetaParse failed: %v", err)
			}
			if !reflect.DeepEqual(output, tt.expected) {
				t.Errorf("Parsed entry doesn't match expected.\nGot: %#v\nWant: %#v", output, tt.expected)
			}
		})
	}
}

func TestMetaParseSyslog_Custom(t *testing.T) {
	module := makeSyslogModule()

	logLine := `Jul 30 14:20:02 myhost sshd[12345]: Accepted password for user from 10.0.0.1 port 22 ssh2`

	expected := map[string]interface{}{
		"timestamp": "Jul 30 14:20:02",
		"host":      "myhost",
		"process":   "sshd",
		"pid":       "12345",
		"message":   "Accepted password for user from 10.0.0.1 port 22 ssh2",
	}

	output, err := MetaParse(logLine, module)
	if err != nil {
		t.Fatalf("MetaParse failed: %v", err)
	}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Parsed entry doesn't match expected.\nGot: %#v\nWant: %#v", output, expected)
	}
}

// -- Apache Tests --

func TestMetaParseApache_Registry(t *testing.T) {
	module := structs.MetaParserRegistry["apache"]

	logLine := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`

	expected := map[string]interface{}{
		"remote":      "127.0.0.1",
		"remote_long": "-",
		"remote_user": "frank",
		"timestamp":   "10/Oct/2000:13:55:36 -0700",
		"request":     "GET /apache_pb.gif HTTP/1.0",
		"status_code": "200",
		"size":        "2326",
	}

	output, err := MetaParse(logLine, module)
	if err != nil {
		t.Fatalf("MetaParse failed: %v", err)
	}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Parsed entry doesn't match expected.\nGot: %#v\nWant: %#v", output, expected)
	}
}

func TestMetaParseApache_Custom(t *testing.T) {
	module := makeApacheModule()

	logLine := `192.168.1.1 - alice [11/Oct/2020:16:22:01 -0700] "POST /login HTTP/1.1" 302 512`

	expected := map[string]interface{}{
		"remote":      "192.168.1.1",
		"remote_long": "-",
		"remote_user": "alice",
		"timestamp":   "11/Oct/2020:16:22:01 -0700",
		"request":     "POST /login HTTP/1.1",
		"status_code": "302",
		"size":        "512",
	}

	output, err := MetaParse(logLine, module)
	if err != nil {
		t.Fatalf("MetaParse failed: %v", err)
	}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Parsed entry doesn't match expected.\nGot: %#v\nWant: %#v", output, expected)
	}
}
