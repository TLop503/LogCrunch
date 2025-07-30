package structs

type Log struct {
	Host      string      `json:"host"`
	Timestamp int64       `json:"timestamp"`
	Type      string      `json:"type"`
	Path      string      `json:"path"`
	Parsed    interface{} `json:"parsed"`
	Raw       string      `json:"raw"`
}

type SyslogEntry struct {
	Timestamp string `json:"timestamp"`
	Host      string `json:"host"`
	Process   string `json:"process"`
	PID       string `json:"pid"`
	Message   string `json:"message"`
}
