package hemoglobin

import (
	"github.com/TLop503/LogCrunch/agent/hemoglobin/modules"
	"github.com/TLop503/LogCrunch/agent/utils"
	"github.com/TLop503/LogCrunch/structs"
	"log"
	"regexp"
	"time"

	"github.com/hpcloud/tail"
)

// ParserEntry holds the regex and a function that returns an empty struct to parse into
type ParserEntry struct {
	Regex     *regexp.Regexp
	NewStruct func() interface{}
}

// MetaParserRegistry maps module names to their regex and struct constructors
var MetaParserRegistry = map[string]ParserEntry{
	"syslog": {
		Regex:     modules.SyslogRegex,
		NewStruct: func() interface{} { return new(structs.SyslogEntry) },
	},
	"apache": {
		Regex:     modules.ApacheRegex,
		NewStruct: func() interface{} { return new(structs.ApacheLogEntry) },
	},
}

// ReadLog watches a log file and parses lines with a generic meta parser
func ReadLog(logChan chan<- structs.Log, config structs.Target) {
	tailConfig := tail.Config{
		ReOpen:    true,
		Follow:    true,
		MustExist: false,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 0},
		Logger:    tail.DiscardingLogger,
	}

	t, err := tail.TailFile(config.Path, tailConfig)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	parserEntry, ok := MetaParserRegistry[config.Module]

	for line := range t.Lines {
		if line.Err != nil {
			log.Printf("Error reading line from file %v: %v\n", config.Path, line.Err)
			continue
		}

		var parsed interface{}

		if ok {
			// Create a new empty struct instance for this line
			outputStruct := parserEntry.NewStruct()

			// Parse line using the generic MetaParse function
			err := modules.MetaParse(line.Text, parserEntry.Regex, outputStruct)
			if err == nil {
				parsed = outputStruct
			} else {
				log.Printf("Parse error for line in %v: %v", config.Path, err)
				parsed = map[string]error{"Parsing error": err}
			}
		}

		logEntry := structs.Log{
			Host:      utils.GetHostName(),
			Timestamp: time.Now().Unix(),
			Type:      config.Name,
			Path:      config.Path,
			Raw:       line.Text,
			Parsed:    parsed,
		}

		logChan <- logEntry
	}
}
