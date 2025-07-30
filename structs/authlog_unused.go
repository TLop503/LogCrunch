package structs

import "time"

type authlog struct {
	Timestamp        time.Time `json:"timestamp"`
	Host             string    `json:"host"`
	Service          string    `json:"service"`
	PID              string    `json:"pid"`
	Event            string    `json:"event"`
	User             string    `json:"user,omitempty"`
	SourceIP         string    `json:"source_ip,omitempty"`
	Port             string    `json:"port,omitempty"`
	Protocol         string    `json:"protocol,omitempty"`
	Command          string    `json:"command,omitempty"`
	TTY              string    `json:"tty,omitempty"`
	WorkingDirectory string    `json:"working_directory,omitempty"`
	TargetUser       string    `json:"target_user,omitempty"`
}
