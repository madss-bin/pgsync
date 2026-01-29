package ui

import (
	"strings"

	"pgsync/internal/db"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleIntro(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.state = StateSourceURL
		m.textInput.Reset()
		m.textInput.Placeholder = "postgres://user:password@host:port/dbname"
		return m, textinput.Blink
	case "h", "H":
		m.state = StateHistory
		return m, loadHistoryCmd()
	}
	return m, nil
}

func (m Model) handleHistory(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = StateIntro
		return m, nil
	}
	return m, nil
}

func (m Model) handleSourceURL(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		url := strings.TrimSpace(m.textInput.Value())
		if err := db.ValidateURL(url); err != nil {
			m.errorMsg = err.Error()
			return m, nil
		}
		m.sourceURL = url
		m.errorMsg = ""
		m.state = StateTargetURL
		m.textInput.Reset()
		m.textInput.Placeholder = "postgres://user:password@host:port/dbname"
		return m, textinput.Blink
	default:
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
}

func (m Model) handleTargetURL(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		url := strings.TrimSpace(m.textInput.Value())
		if err := db.ValidateURL(url); err != nil {
			m.errorMsg = err.Error()
			return m, nil
		}
		if err := db.URLsAreDifferent(m.sourceURL, url); err != nil {
			m.errorMsg = err.Error()
			return m, nil
		}
		m.targetURL = url
		m.errorMsg = ""

		m.state = StateEstimation
		return m, estimateCmd(m.sourceURL, m.targetURL)
	default:
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
}

func (m Model) handleEstimation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.state = StateTableSelect
		m.cursor = 0
		return m, fetchTablesCmd(m.sourceURL)
	}
	return m, nil
}

func (m Model) handleTableSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			if m.cursor < m.scrollOffset {
				m.scrollOffset = m.cursor
			}
		}
	case "down", "j":
		if m.cursor < len(m.availableTables)-1 {
			m.cursor++
			if m.cursor >= m.scrollOffset+10 {
				m.scrollOffset = m.cursor - 9
			}
		}
	case " ":
		if len(m.availableTables) > 0 {
			table := m.availableTables[m.cursor]
			if m.selectedTables[table] {
				delete(m.selectedTables, table)
			} else {
				m.selectedTables[table] = true
			}
		}
	case "a", "A":
		if len(m.selectedTables) == len(m.availableTables) {
			m.selectedTables = make(map[string]bool)
		} else {
			for _, t := range m.availableTables {
				m.selectedTables[t] = true
			}
		}
	case "enter":
		m.options.SelectedTables = []string{}
		for t := range m.selectedTables {
			m.options.SelectedTables = append(m.options.SelectedTables, t)
		}

		m.state = StateOptions
		m.cursor = 0
		return m, nil
	}
	return m, nil
}

func (m Model) handleOptions(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < 2 {
			m.cursor++
		}
	case "right", "l":
		if m.cursor == 0 {
			if m.options.ParallelJobs < 16 {
				m.options.ParallelJobs++
			}
		}
	case "left", "h":
		if m.cursor == 0 {
			if m.options.ParallelJobs > 1 {
				m.options.ParallelJobs--
			}
		}
	case " ", "enter":
		switch m.cursor {
		case 1:
			m.options.AutoBackup = !m.options.AutoBackup
		case 2:
			m.state = StateMigrationType
			m.selectedIndex = 0
			return m, nil
		}
	}
	return m, nil
}

func (m Model) handleMigrationType(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.selectedIndex > 0 {
			m.selectedIndex--
		}
	case "down", "j":
		if m.selectedIndex < 2 {
			m.selectedIndex++
		}
	case "enter":
		switch m.selectedIndex {
		case 0:
			m.migrationType = db.SchemaAndData
		case 1:
			m.migrationType = db.SchemaOnly
		case 2:
			m.migrationType = db.DataOnly
		}
		m.state = StateMigrating
		m.progressChan = make(chan db.ProgressUpdate, 100)
		return m, m.startMigration()
	}
	return m, nil
}

func (m Model) startMigration() tea.Cmd {
	progressChan := m.progressChan
	src := m.sourceURL
	tgt := m.targetURL
	typ := m.migrationType
	opts := m.options

	go func() {
		migrator := db.NewMigrator(src, tgt, typ, opts, progressChan)
		stats, err := migrator.Migrate()
		if err != nil {
			progressChan <- db.ProgressUpdate{
				Percentage: -1,
				Message:    err.Error(),
				Stats:      stats,
			}
		} else {
			progressChan <- db.ProgressUpdate{
				Percentage: 1.0,
				Message:    "Migration completed!",
				Stats:      stats,
			}
		}
		close(progressChan)
	}()

	return tea.Batch(
		waitForMigrationUpdate(progressChan),
		tickCmd(),
	)
}
