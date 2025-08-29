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

/*
	Hemoglobin is a <routine> containing <data> that facilitates the transportation of <logs> in <agents>.
*/

// ParserEntry holds the regex and a function that returns an empty struct to parse into
type ParserEntry struct {
	Regex     *regexp.Regexp
	NewStruct func() interface{}
}

// ReadLog watches a log file and parses lines with a generic meta parser
func ReadLog(logChan chan<- structs.Log, target structs.Target) {
	parserModule, err := modules.HandleConfigTarget(target)
	if err != nil {
		log.Println("Error handling config target:", err)
		return
	}

	tailConfig := tail.Config{
		ReOpen:    true,                                 // handle rotation
		Follow:    true,                                 // continuous
		MustExist: false,                                // don't err if file dne yet
		Location:  &tail.SeekInfo{Offset: 0, Whence: 0}, // start from end of file
		Logger:    tail.DiscardingLogger,                // disable internal logging
	}

	t, err := tail.TailFile(target.Path, tailConfig)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	for line := range t.Lines {
		if line.Err != nil {
			log.Printf("Error reading line from file %v: %v\n", target.Path, line.Err)
			continue
		}

		var parsed interface{}

		// Parse line using the generic MetaParse function
		parsed, err := modules.MetaParse(line.Text, parserModule)
		if err != nil {
			log.Printf("Parse error for line in %v: %v", target.Path, err)
			parsed = map[string]error{"Parsing error": err}
		}

		logEntry := structs.Log{
			Host:      utils.GetHostName(),
			Timestamp: time.Now().Unix(),
			Type:      target.Module,
			Name:      target.Name,
			Path:      target.Path,
			Raw:       line.Text,
			Parsed:    parsed,
		}

		logChan <- logEntry
	}
}
