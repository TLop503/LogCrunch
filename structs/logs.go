package structs

type Log struct {
	Host      string      `json:"host"`
	Timestamp int64       `json:"timestamp"`
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	Path      string      `json:"path"`
	Parsed    interface{} `json:"parsed"`
	Raw       string      `json:"raw"`
}

type SyslogEntry struct {
	Timestamp string `logfield:"timestamp"`
	Host      string `logfield:"host"`
	Process   string `logfield:"process"`
	PID       string `logfield:"pid"`
	Message   string `logfield:"message"`
}

type ApacheLogEntry struct {
	Remote     string `logfield:"remote"`
	RemoteLong string `logfield:"remote_long"`
	RemoteUser string `logfield:"remote_user"`
	Timestamp  string `logfield:"timestamp"`
	Request    string `logfield:"request"`
	StatusCode string `logfield:"status_code"`
	Size       string `logfield:"size"`
}

type NginxLogEntry struct {
	ApacheFormattedData string `json:"apacheFormattedData"`
	Referrer            string `json:"referrer"`
	UserAgent           string `json:"userAgent"`
}
