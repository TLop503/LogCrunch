package structs

type Log struct {
	Host      string      `json:"host"`
	Timestamp int64       `json:"timestamp"`
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
}

type IntakeLogFileData struct {
	FileContent string `json:"fileContent"`
}
