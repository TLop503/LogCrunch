package structs

import (
	"reflect"
	"regexp"
)

// MetaParserRegistry maps module names to their regex and struct constructors
var MetaParserRegistry = map[string]ParserModule{
	"syslog": {
		Regex:  regexp.MustCompile(`^(?P<timestamp>\w+\s+\d+\s+\d+:\d+:\d+)\s+(?P<host>\S+)\s+(?P<process>\w+)(?:\[(?P<pid>\d+)\])?:\s+(?P<message>.*)$`),
		Schema: ReflectSchema(SyslogEntry{}),
	},
	"apache": {
		Regex:  regexp.MustCompile(`(?P<remote>\S+) (?P<remote_long>\S+) (?P<remote_user>\S+) \[(?P<timestamp>[^\]]+)\] "(?P<request>[^"]*)" (?P<status_code>\d{3}) (?P<size>\S+)`),
		Schema: ReflectSchema(ApacheLogEntry{}),
	},
	"Heartbeat": {
		Regex:  regexp.MustCompile(`/ \d+ /`),
		Schema: map[string]string{},
	},
}

// ReflectSchema takes a struct instance and returns a map of field name -> type string
// It reads the `logfield` tag if present, otherwise uses the Go field name.
func ReflectSchema(s interface{}) map[string]string {
	schema := make(map[string]string)

	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Field name: prefer `logfield` tag
		name := field.Tag.Get("logfield")
		if name == "" {
			name = field.Name
		}

		// Field type: convert Go type to simple string representation
		var typeStr string
		switch field.Type.Kind() {
		case reflect.String:
			typeStr = "string"
		case reflect.Int, reflect.Int64, reflect.Int32:
			typeStr = "int"
		case reflect.Float32, reflect.Float64:
			typeStr = "float"
		case reflect.Bool:
			typeStr = "bool"
		default:
			typeStr = "interface{}"
		}

		schema[name] = typeStr
	}

	return schema
}
