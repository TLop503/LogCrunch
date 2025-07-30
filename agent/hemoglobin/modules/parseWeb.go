package modules

import (
	"errors"
	"regexp"

	"github.com/TLop503/LogCrunch/structs"
)

var apacheRegex = regexp.MustCompile(`(?P<remote>\S+) (?P<remote_long>\S+) (?P<remote_user>\S+) \[(?P<timestamp>[^\]]+)\] "(?P<request>[^"]*)" (?P<status_code>\d{3}) (?P<size>\S+)`)

func ParseApacheLog(line string) (*structs.ApacheLogEntry, error) {
	matches := apacheRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, errors.New("failed to parse apache log line")
	}

	entry := &structs.ApacheLogEntry{}
	for i, name := range apacheRegex.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		switch name {
		case "remote":
			entry.Remote = matches[i]
		case "remote_long":
			entry.RemoteLong = matches[i]
		case "remote_user":
			entry.RemoteUser = matches[i]
		case "timestamp":
			entry.Timestamp = matches[i]
		case "request":
			entry.Request = matches[i]
		case "status_code":
			entry.StatusCode = matches[i]
		case "size":
			entry.Size = matches[i]
		}
	}

	return entry, nil
}
