package modules

import (
	"fmt"
	"github.com/TLop503/LogCrunch/structs"
	"regexp"
)

var syslogRegex = regexp.MustCompile(`^(?P<timestamp>\w+\s+\d+\s+\d+:\d+:\d+)\s+(?P<hostname>\S+)\s+(?P<process>\w+)(?:\[(?P<pid>\d+)\])?:\s+(?P<message>.*)$`)

func ParseSyslog(line string) (*structs.SyslogEntry, error) {
	match := syslogRegex.FindStringSubmatch(line)
	if match == nil {
		return nil, fmt.Errorf("no match")
	}

	result := &structs.SyslogEntry{}
	for i, name := range syslogRegex.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		switch name {
		case "timestamp":
			result.Timestamp = match[i]
		case "hostname":
			result.Host = match[i]
		case "process":
			result.Process = match[i]
		case "pid":
			result.PID = match[i]
		case "message":
			result.Message = match[i]
		}
	}
	return result, nil
}
