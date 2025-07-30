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

type Target struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	Severity string `yaml:"severity"`
	Module   string `yaml:"module"`
}

type YamlConfig struct {
	Targets []Target `yaml:"Targets"`
}
