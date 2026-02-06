package structs

import "regexp"

type IntakeLogFileData struct {
	FileContent string `json:"fileContent"`
}

type Target struct {
	Name     string            `yaml:"name"`
	Path     string            `yaml:"path"`
	Severity string            `yaml:"severity"`
	Custom   bool              `yaml:"custom"`
	Module   string            `yaml:"module,omitempty"`
	Regex    string            `yaml:"regex,omitempty"`
	Schema   map[string]string `yaml:"schema,omitempty"`
}

type Service struct {
	Name     string `yaml:"name"`
	Key      string `yaml:"key"`
	Severity string `yaml:"severity"`
}

type YamlConfig struct {
	Targets  []Target  `yaml:"Targets"`
	Services []Service `yaml:"Services"`
}

type ParserModule struct {
	Regex  *regexp.Regexp
	Schema map[string]string // field name -> type
}
