package hemoglobin

import (
	"github.com/TLop503/LogCrunch/agent/hemoglobin/modules"
	"github.com/TLop503/LogCrunch/agent/utils"
	"github.com/TLop503/LogCrunch/structs"
	"log"
	"time"

	"github.com/hpcloud/tail"
)

/*
Hemoglobin is a <routine> containing <data> that facilitates the transportation of <logs> in <agents>.
*/

// wrap() is a helper function to prepare parsing modules for inclusion in the registry
func wrap[T any](fn func(string) (*T, error)) func(string) (interface{}, error) {
	return func(line string) (interface{}, error) {
		return fn(line)
	}
}

// mapping of modules specifiable in the config and written parsing modules
var ParserRegistry = map[string]func(string) (interface{}, error){
	"syslog": wrap(modules.ParseSyslog),
}

// watch a log file for updates as they come in
func ReadLog(logChan chan<- structs.Log, config structs.Target) {
	var seekOffset int64 = 0

	tailConfig := tail.Config{
		ReOpen:    true,                                          // Handle log rotation
		Follow:    true,                                          // Continuously read new lines
		MustExist: false,                                         // Don't error if the file doesn't exist initially
		Location:  &tail.SeekInfo{Offset: seekOffset, Whence: 0}, // Start from the given offset (0 = beginning)
		Logger:    tail.DiscardingLogger,                         // Disable internal logging
	}

	t, err := tail.TailFile(config.Path, tailConfig)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	parser := ParserRegistry[config.Module]

	for line := range t.Lines {
		if line.Err != nil {
			log.Printf("Error reading line from file %v: %v\n", config.Path, line.Err)
			continue
		}

		var parsed interface{}
		if parser != nil {
			parsed, _ = parser(line.Text) // ignore error for now
		}

		logEntry := structs.Log{
			Host:      utils.GetHostName(),
			Timestamp: time.Now().Unix(),
			Type:      config.Name,
			Path:      config.Path,
			Raw:       line.Text,
			Parsed:    parsed,
		}

		// send to channel for writing across wire
		logChan <- logEntry
	}
}
