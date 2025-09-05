package modules

import (
	"fmt"
	"strconv"

	"github.com/TLop503/LogCrunch/structs"
)

// MetaParse parses logs into named fields and stages them for transmission
// Takes a log (read from a file), and a module containing a regex and a schema
// The log is processed with the regex, and the resulting named fields are
// organized according to the schema
// the final logcrunch log is sent to the log channel for transmission
// to the siem server
func MetaParse(log string, module structs.ParserModule) (map[string]interface{}, error) {
	match := module.Regex.FindStringSubmatch(log)
	if match == nil {
		return nil, fmt.Errorf("no match found")
	}

	names := module.Regex.SubexpNames()
	if len(names) != len(match) {
		return nil, fmt.Errorf("capture group count mismatch")
	}

	parsedLog := make(map[string]interface{})

	for i, name := range names {
		if i == 0 || name == "" {
			continue // skip full match or unnamed groups
		}

		// Determine type from schema
		fieldType, ok := module.Schema[name]
		if !ok {
			// if the schema doesn't include this field, just store as string
			parsedLog[name] = match[i]
			continue
		}

		// attempt parse numbers to correct type
		// if conversion fails, assign as strings instead
		switch fieldType {
		case "int":
			val, err := strconv.Atoi(match[i])
			if err != nil {
				parsedLog[name] = match[i]
			}
			parsedLog[name] = val
		case "float":
			val, err := strconv.ParseFloat(match[i], 64)
			if err != nil {
				parsedLog[name] = match[i]
			}
			parsedLog[name] = val
		default:
			// if not parsed to number just use string
			parsedLog[name] = match[i]
		}
	}

	return parsedLog, nil
}
