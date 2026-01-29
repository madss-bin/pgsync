package db

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type MigrationStatus string

const (
	StatusSuccess MigrationStatus = "success"
	StatusFailed  MigrationStatus = "failed"
)

type MigrationRecord struct {
	Timestamp     time.Time       `json:"timestamp"`
	Source        string          `json:"source"`
	Target        string          `json:"target"`
	MigrationType MigrationType   `json:"migration_type"`
	Status        MigrationStatus `json:"status"`
	Duration      string          `json:"duration"`
	Error         string          `json:"error,omitempty"`
}

func getHistoryPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".pgsync")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "history.json"), nil
}

func SaveHistory(record MigrationRecord) error {
	path, err := getHistoryPath()
	if err != nil {
		return err
	}

	var history []MigrationRecord
	data, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(data, &history)
	}

	history = append([]MigrationRecord{record}, history...)

	if len(history) > 50 {
		history = history[:50]
	}

	newData, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, newData, 0644)
}

func LoadHistory() ([]MigrationRecord, error) {
	path, err := getHistoryPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []MigrationRecord{}, nil
		}
		return nil, err
	}

	var history []MigrationRecord
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}

	sort.Slice(history, func(i, j int) bool {
		return history[i].Timestamp.After(history[j].Timestamp)
	})

	return history, nil
}
