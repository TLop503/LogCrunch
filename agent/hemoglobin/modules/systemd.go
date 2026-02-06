package modules

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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

// ListenToSystemd Logs through coreos journal
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
			continue
		}
		parsed, err := entryToPrettyString(entry)
		if err != nil {
			log.Printf("Error parsing journal entry: %s", err)
		}
		raw, err := entryToString(entry)
		if err != nil {
			log.Printf("Error converting journal entry to string: %s", err)
		}
		logEntry := structs.Log{
			Host:      utils.GetHostName(),
			Timestamp: time.Now().Unix(),
			Module:    "systemd",
			Name:      entry.Fields["_SYSTEMD_UNIT"],
			Path:      "systemd",
			Raw:       raw,
			Parsed:    parsed,
		}
		logChan <- logEntry
	}
}

// systemd entryToString
func entryToString(entry *sdjournal.JournalEntry) (string, error) {
	// Serialize the Fields map to JSON
	out, err := json.MarshalIndent(entry.Fields, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal journal entry: %w", err)
	}
	return string(out), nil
}

// entryToPrettyString creates "parsed" entries
func entryToPrettyString(entry *sdjournal.JournalEntry) (structs.SyslogPrettyEntry, error) {
	priorityInt, err := strconv.Atoi(entry.Fields["PRIORITY"])
	if err != nil {
		priorityInt = -1
	}

	out := structs.SyslogPrettyEntry{
		Message:  entry.Fields["MESSAGE"],
		Priority: priorityInt,
		Cmdline:  entry.Fields["_CMDLINE"],
	}

	return out, nil
}
