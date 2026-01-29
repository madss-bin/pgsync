package db

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type MigrationType string

const (
	SchemaAndData MigrationType = "schema_data"
	SchemaOnly    MigrationType = "schema_only"
	DataOnly      MigrationType = "data_only"
)

type ProgressUpdate struct {
	Percentage float64
	Message    string
	Command    string
	Stats      *MigrationStats
}

type MigrationOptions struct {
	SelectedTables []string
	ParallelJobs   int
	AutoBackup     bool
}

type MigrationStats struct {
	MigrationType   MigrationType
	TablesMigrated  int
	Duration        string
	Warnings        []string
	BackupPath      string
	DidRollback     bool
	RollbackSuccess bool
	LogPath         string
}

type Migrator struct {
	source        string
	target        string
	migrationType MigrationType
	options       MigrationOptions
	progressChan  chan<- ProgressUpdate
	stats         MigrationStats
	backupPath    string
}

func NewMigrator(source, target string, migrationType MigrationType, options MigrationOptions, progressChan chan<- ProgressUpdate) *Migrator {
	return &Migrator{
		source:        source,
		target:        target,
		migrationType: migrationType,
		options:       options,
		progressChan:  progressChan,
	}
}

func (m *Migrator) Migrate() (*MigrationStats, error) {
	startTime := time.Now()
	var finalErr error

	m.stats = MigrationStats{
		MigrationType: m.migrationType,
		Warnings:      []string{},
	}

	logFile, err := os.CreateTemp("", "pgsync_migration_*.log")
	if err == nil {
		m.stats.LogPath = logFile.Name()
		defer logFile.Close()
	}

	writeLog := func(format string, args ...interface{}) {
		if logFile != nil {
			msg := fmt.Sprintf(format, args...)
			timestamp := time.Now().Format("15:04:05")
			logFile.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, msg))
		}
	}

	writeLog("Starting migration %s -> %s (Type: %s)", redactURL(m.source), redactURL(m.target), m.migrationType)

	if len(m.options.SelectedTables) > 0 {
		m.stats.TablesMigrated = len(m.options.SelectedTables)
	}

	defer func() {
		m.stats.Duration = time.Since(startTime).Round(time.Second).String()

		status := StatusSuccess
		errMsg := ""
		if finalErr != nil {
			status = StatusFailed
			errMsg = finalErr.Error()
		}

		record := MigrationRecord{
			Timestamp:     startTime,
			Source:        redactURL(m.source),
			Target:        redactURL(m.target),
			MigrationType: m.migrationType,
			Status:        status,
			Duration:      m.stats.Duration,
			Error:         errMsg,
		}
		_ = SaveHistory(record)
	}()

	m.sendProgress(0.0, "Preparing migration tasks...", "")

	m.sendProgress(0.1, "Step 1/5: Verifying source connection...", "pg_isready -d "+redactURL(m.source))
	if err := CheckConnection(m.source); err != nil {
		finalErr = fmt.Errorf("source database: %w", err)
		return &m.stats, finalErr
	}

	m.sendProgress(0.2, "Step 1/5: Verifying target connection...", "pg_isready -d "+redactURL(m.target))
	if err := CheckConnection(m.target); err != nil {
		finalErr = fmt.Errorf("target database: %w", err)
		return &m.stats, finalErr
	}

	if m.options.AutoBackup {
		backupFile := fmt.Sprintf("backup_target_%d.dump", time.Now().Unix())
		m.sendProgress(0.3, "Step 2/5: Creating safety backup of target...", "pg_dump ... -w > "+backupFile)

		backupCmd := exec.Command("pg_dump", "-d", m.target, "-w", "-Fc", "-f", backupFile)
		writeLog("Running backup command: %v", backupCmd.Args)
		if out, err := backupCmd.CombinedOutput(); err != nil {
			writeLog("Backup failed: %s", string(out))
			warning := fmt.Sprintf("Safety backup failed: %s", string(out))
			m.stats.Warnings = append(m.stats.Warnings, warning)
			m.sendProgress(0.3, "Warning: Safety backup failed, proceeding...", string(out))
		} else {
			m.backupPath = backupFile
			m.stats.BackupPath = backupFile
			m.sendProgress(0.35, "Step 2/5: Safety backup created: "+backupFile, "")
		}
	}

	m.sendProgress(0.4, "Step 3/5: Dumping source database...", "")
	tmpFile, err := os.CreateTemp("", "pgsync_dump_*.dump")
	if err != nil {
		finalErr = fmt.Errorf("failed to create temp file: %w", err)
		return &m.stats, finalErr
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	args := []string{m.source, "-w", "-Fc", "-f", tmpPath}
	switch m.migrationType {
	case SchemaOnly:
		args = append(args, "--schema-only")
	case DataOnly:
		args = append(args, "--data-only")
	}

	if len(m.options.SelectedTables) > 0 {
		for _, t := range m.options.SelectedTables {
			args = append(args, "-t", t)
		}
	}

	args = append(args, "--no-owner", "--no-privileges", "--verbose")

	dumpCmdStr := fmt.Sprintf("pg_dump %s ... (tables: %d)", redactURL(m.source), len(m.options.SelectedTables))
	m.sendProgress(0.5, "Step 3/5: Dumping content...", dumpCmdStr)

	dumpCmd := exec.Command("pg_dump", args...)
	writeLog("Running dump command: %v", dumpCmd.Args)
	if output, err := dumpCmd.CombinedOutput(); err != nil {
		writeLog("Dump failed: %s", string(output))
		finalErr = fmt.Errorf("dump failed: %s", string(output))
		return &m.stats, finalErr
	}

	if fi, err := os.Stat(tmpPath); err == nil {
		writeLog("Dump file size: %d bytes", fi.Size())
		if fi.Size() == 0 {
			finalErr = fmt.Errorf("dump file is empty (0 bytes) - check source database permissions or connectivity")
			writeLog("Error: Dump file is empty")
			return &m.stats, finalErr
		}
	} else {
		writeLog("Warning: Could not check dump file size: %v", err)
	}

	jobs := "4"
	if m.options.ParallelJobs > 0 {
		jobs = fmt.Sprintf("%d", m.options.ParallelJobs)
	}

	restoreArgs := []string{"-d", m.target, "-w", "-j", jobs, "-c", "--if-exists", "--no-owner", "--no-privileges", "--verbose", tmpPath}
	restoreCmdStr := fmt.Sprintf("pg_restore -d %s -j %s ...", redactURL(m.target), jobs)

	m.sendProgress(0.8, fmt.Sprintf("Step 4/5: Parallel restore (j=%s)...", jobs), restoreCmdStr)

	restoreCmd := exec.Command("pg_restore", restoreArgs...)
	writeLog("Running restore command: %v", restoreCmd.Args)
	output, err := restoreCmd.CombinedOutput()

	if err != nil {
		outputStr := string(output)
		if strings.Contains(outputStr, `unrecognized configuration parameter "transaction_timeout"`) {
			writeLog("Restore returned error matches version mismatch pattern, treating as warning: %v", err)
			writeLog("Output: %s", outputStr)
			m.stats.Warnings = append(m.stats.Warnings, "Ignored benign 'transaction_timeout' errors (PG 17 -> Older DB)")
		} else {
			writeLog("Restore failed: %s", outputStr)
			restoreErr := fmt.Errorf("restore failed: %s", outputStr)

			if m.backupPath != "" {
				m.stats.DidRollback = true
				m.sendProgress(0.85, "Restore failed! Attempting rollback from backup...", "")

				rollbackArgs := []string{"-d", m.target, "-w", "-c", "--if-exists", "-j", jobs, m.backupPath}
				rollbackCmd := exec.Command("pg_restore", rollbackArgs...)
				writeLog("Running rollback command: %v", rollbackCmd.Args)

				if rbOut, rbErr := rollbackCmd.CombinedOutput(); rbErr != nil {
					writeLog("Rollback failed: %s", string(rbOut))
					m.stats.RollbackSuccess = false
					m.stats.Warnings = append(m.stats.Warnings, fmt.Sprintf("Rollback also failed: %s", string(rbOut)))
					finalErr = fmt.Errorf("%v (rollback also failed)", restoreErr)
				} else {
					writeLog("Rollback successful")
					m.stats.RollbackSuccess = true
					m.sendProgress(0.9, "Rollback successful! Target database restored to previous state.", "")
					finalErr = fmt.Errorf("%v (rolled back successfully)", restoreErr)
				}
			} else {
				finalErr = restoreErr
			}

			return &m.stats, finalErr
		}
	}

	m.sendProgress(1.0, "Step 5/5: Migration completed!", "")

	return &m.stats, nil
}

func (m *Migrator) sendProgress(percentage float64, message string, command string) {
	if m.progressChan != nil {
		m.progressChan <- ProgressUpdate{
			Percentage: percentage,
			Message:    message,
			Command:    command,
		}
	}
}

func redactURL(url string) string {
	return url
}
