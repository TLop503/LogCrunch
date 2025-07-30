package modules

import (
	"testing"

	"github.com/TLop503/LogCrunch/structs"
)

func TestParseSyslog(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *structs.SyslogEntry
		wantErr  bool
	}{
		{
			name:  "basic syslog entry",
			input: "Jul 29 12:34:56 myhost systemd[1]: Started Session 1 of user root.",
			expected: &structs.SyslogEntry{
				Timestamp: "Jul 29 12:34:56",
				Host:      "myhost",
				Process:   "systemd",
				PID:       "1",
				Message:   "Started Session 1 of user root.",
			},
			wantErr: false,
		},
		{
			name:  "entry without PID",
			input: "Jul 29 12:34:56 myhost cron: Job started.",
			expected: &structs.SyslogEntry{
				Timestamp: "Jul 29 12:34:56",
				Host:      "myhost",
				Process:   "cron",
				PID:       "",
				Message:   "Job started.",
			},
			wantErr: false,
		},
		{
			name:    "malformed entry",
			input:   "this is not a valid syslog line",
			wantErr: true,
		},
		{
			name:  "process name with digits",
			input: "Jul 29 12:34:56 localhost kernel123[999]: A strange message",
			expected: &structs.SyslogEntry{
				Timestamp: "Jul 29 12:34:56",
				Host:      "localhost",
				Process:   "kernel123",
				PID:       "999",
				Message:   "A strange message",
			},
			wantErr: false,
		},
		{
			name:  "multispace message",
			input: "Jul 29 12:34:56 host1 sshd[2222]: Accepted password for user from 192.168.0.1 port 22 ssh2",
			expected: &structs.SyslogEntry{
				Timestamp: "Jul 29 12:34:56",
				Host:      "host1",
				Process:   "sshd",
				PID:       "2222",
				Message:   "Accepted password for user from 192.168.0.1 port 22 ssh2",
			},
			wantErr: false,
		},
		{
			name:  "REAL cron with pam_unix message",
			input: "Jul 30 14:17:01 blackwall CRON[620010]: pam_unix(cron:session): session opened for user root(uid=0) by (uid=0)",
			expected: &structs.SyslogEntry{
				Timestamp: "Jul 30 14:17:01",
				Host:      "blackwall",
				Process:   "CRON",
				PID:       "620010",
				Message:   "pam_unix(cron:session): session opened for user root(uid=0) by (uid=0)",
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseSyslog(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if got == nil {
				t.Errorf("got nil result")
				return
			}

			if *got != *tc.expected {
				t.Errorf("unexpected result:\ngot:  %+v\nwant: %+v", got, tc.expected)
			}
		})
	}
}
