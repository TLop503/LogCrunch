package modules

import (
	"reflect"
	"testing"

	"github.com/TLop503/LogCrunch/structs"
)

// Helper to create a sample custom target
func makeCustomTarget() structs.Target {
	return structs.Target{
		Name:   "CustomSyslog",
		Path:   "/var/log/custom.log",
		Custom: true,
		Regex:  `^(?P<timestamp>\w+\s+\d+\s+\d+:\d+:\d+)\s+(?P<host>\S+)\s+(?P<process>\w+)(?:\[(?P<pid>\d+)\])?:\s+(?P<message>.*)$`,
		Schema: map[string]string{
			"timestamp": "string",
			"host":      "string",
			"process":   "string",
			"pid":       "string",
			"message":   "string",
		},
	}
}

// -- Tests --

func TestHandleConfigTarget_Custom(t *testing.T) {
	target := makeCustomTarget()
	module, err := HandleConfigTarget(target)
	if err != nil {
		t.Fatalf("HandleConfigTarget failed: %v", err)
	}

	// Regex should match the one in the target
	if module.Regex.String() != target.Regex {
		t.Errorf("Expected regex %s, got %s", target.Regex, module.Regex.String())
	}

	// Schema should be the same as the target
	if !reflect.DeepEqual(module.Schema, target.Schema) {
		t.Errorf("Expected schema %#v, got %#v", target.Schema, module.Schema)
	}
}

func TestHandleConfigTarget_RegistrySyslog(t *testing.T) {
	target := structs.Target{
		Name:   "SyslogTarget",
		Path:   "/var/log/auth.log",
		Custom: false,
		Module: "syslog",
	}

	module, err := HandleConfigTarget(target)
	if err != nil {
		t.Fatalf("HandleConfigTarget failed: %v", err)
	}

	registryModule := structs.MetaParserRegistry["syslog"]

	if module.Regex.String() != registryModule.Regex.String() {
		t.Errorf("Expected regex %s, got %s", registryModule.Regex.String(), module.Regex.String())
	}

	if !reflect.DeepEqual(module.Schema, registryModule.Schema) {
		t.Errorf("Expected schema %#v, got %#v", registryModule.Schema, module.Schema)
	}
}

func TestHandleConfigTarget_RegistryApache(t *testing.T) {
	target := structs.Target{
		Name:   "ApacheTarget",
		Path:   "/var/log/apache.log",
		Custom: false,
		Module: "apache",
	}

	module, err := HandleConfigTarget(target)
	if err != nil {
		t.Fatalf("HandleConfigTarget failed: %v", err)
	}

	registryModule := structs.MetaParserRegistry["apache"]

	if module.Regex.String() != registryModule.Regex.String() {
		t.Errorf("Expected regex %s, got %s", registryModule.Regex.String(), module.Regex.String())
	}

	if !reflect.DeepEqual(module.Schema, registryModule.Schema) {
		t.Errorf("Expected schema %#v, got %#v", registryModule.Schema, module.Schema)
	}
}

func TestHandleConfigTarget_InvalidRegistryModule(t *testing.T) {
	target := structs.Target{
		Name:   "InvalidTarget",
		Path:   "/var/log/invalid.log",
		Custom: false,
		Module: "nonexistent",
	}

	_, err := HandleConfigTarget(target)
	if err == nil {
		t.Fatal("Expected error for invalid registry module, got nil")
	}
}

func TestHandleConfigTarget_InvalidRegex(t *testing.T) {
	target := structs.Target{
		Name:   "BadRegex",
		Path:   "/var/log/bad.log",
		Custom: true,
		Regex:  `^(unclosed[`,
		Schema: map[string]string{"field": "string"},
	}

	_, err := HandleConfigTarget(target)
	if err == nil {
		t.Fatal("Expected error for invalid regex, got nil")
	}
}
