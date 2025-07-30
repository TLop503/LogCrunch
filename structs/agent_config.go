package structs

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
