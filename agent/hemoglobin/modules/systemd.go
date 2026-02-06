package modules

import (
	"fmt"
	"log"
	"time"

	"github.com/TLop503/LogCrunch/agent/utils"
	"github.com/TLop503/LogCrunch/structs"
	"github.com/coreos/go-systemd/v22/sdjournal"
)

// StartJour creates a new systemd journal reader
func startJour() (*sdjournal.Journal, error) {
	j, err := sdjournal.NewJournal()
	if err != nil {
		return nil, fmt.Errorf("Error starting journal: %s", err)
	}
	return j, nil
}

func ListenToSystemd(logChan chan<- structs.Log, services []structs.Service) {
	j, err := startJour()

	if err != nil {
		log.Fatalf("Error starting jour: %s", err)
	}

	defer j.Close()

	// add services to listener
	for _, service := range services {
		// parse keys (daemon names, usually) into services
		err := j.AddMatch("_SYSTEMD_UNIT=" + service.Key + ".service")
		if err != nil {
			log.Printf("Error adding systemd unit: %s", err)
		}

	}

	// move cursor to now
	j.SeekHead()
	for {
		n, err := j.Next()
		if err != nil {
			log.Fatalf("Error reading systemd journal: %s", err)
		}
		if n == 0 {
			break
		}

		entry, err := j.GetEntry()
		if err != nil {
			log.Printf("Error getting journal entry: %s", err)
		}
		logEntry := structs.Log{
			Host:      utils.GetHostName(),
			Timestamp: time.Now().Unix(),
			Module:    "systemd",
			Name:      entry.Fields["SYSLOG_IDENTIFIER"],
			Path:      "systemd",
			Raw:       entry.Fields["MESSAGE"],
			Parsed:    entry.Fields["MESSAGE"],
		}
		logChan <- logEntry
	}
}
