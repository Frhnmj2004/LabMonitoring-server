package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"

	"github.com/Frhnmj2004/LabMonitoring-server/models"
)

const bufferFile = "buffer.log"

var (
	bufferMu sync.Mutex
)

// WriteBuffer writes a resource log to the buffer file when offline
func WriteBuffer(log *models.ResourceLog) error {
	bufferMu.Lock()
	defer bufferMu.Unlock()

	file, err := os.OpenFile(bufferFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.Marshal(log)
	if err != nil {
		return err
	}

	if _, err := file.Write(append(jsonData, '\n')); err != nil {
		return err
	}

	return nil
}

// SyncBuffer reads and processes buffered logs, then clears the buffer
func SyncBuffer(processor func(log *models.ResourceLog) error) error {
	bufferMu.Lock()
	defer bufferMu.Unlock()

	// Check if buffer file exists
	if _, err := os.Stat(bufferFile); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(bufferFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var logs []*models.ResourceLog

	// Read all logs from buffer
	for scanner.Scan() {
		var log models.ResourceLog
		if err := json.Unmarshal(scanner.Bytes(), &log); err != nil {
			continue // Skip invalid entries
		}
		logs = append(logs, &log)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Process all logs
	for _, log := range logs {
		if err := processor(log); err != nil {
			return err
		}
	}

	// Clear buffer file after successful processing
	return os.Truncate(bufferFile, 0)
}
