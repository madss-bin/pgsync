package ui

import (
	"time"

	"pgsync/internal/db"
	"pgsync/internal/pkgmgr"

	tea "github.com/charmbracelet/bubbletea"
)

type DepsCheckedMsg struct{ Missing []string }
type DepsInstalledMsg struct{ Err error }
type ProgressMsg struct {
	Percentage float64
	Message    string
	Command    string
	Stats      *db.MigrationStats
}
type MigrationCompleteMsg struct{}
type MigrationErrorMsg string
type TickMsg time.Time

type EstimationMsg struct {
	Result *db.EstimationResult
	Err    error
}

type TablesMsg struct {
	Tables []string
	Err    error
}

type HistoryMsg struct {
	History []db.MigrationRecord
	Err     error
}

func checkDepsCmd() tea.Msg {
	err := db.CheckDependencies()
	var missing []string
	if err != nil {
		missing = append(missing, err.Error())
	}
	return DepsCheckedMsg{Missing: missing}
}

func installDepsCmd() tea.Msg {
	err := pkgmgr.InstallPostgreSQL()
	return DepsInstalledMsg{Err: err}
}

func waitForMigrationUpdate(progressChan chan db.ProgressUpdate) tea.Cmd {
	return func() tea.Msg {
		update, ok := <-progressChan
		if !ok {
			return MigrationCompleteMsg{}
		}
		return ProgressMsg{
			Percentage: update.Percentage,
			Message:    update.Message,
			Command:    update.Command,
			Stats:      update.Stats,
		}
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func estimateCmd(source, target string) tea.Cmd {
	return func() tea.Msg {
		res, err := db.Estimate(source, target)
		return EstimationMsg{Result: res, Err: err}
	}
}

func fetchTablesCmd(url string) tea.Cmd {
	return func() tea.Msg {
		tables, err := db.GetTables(url)
		return TablesMsg{Tables: tables, Err: err}
	}
}

func loadHistoryCmd() tea.Cmd {
	return func() tea.Msg {
		hist, err := db.LoadHistory()
		return HistoryMsg{History: hist, Err: err}
	}
}
