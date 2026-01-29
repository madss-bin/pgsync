package ui

import (
	"pgsync/internal/db"
	"pgsync/internal/pkgmgr"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type State int

const (
	StateCheckingDeps State = iota
	StateInstallingDeps
	StateIntro
	StateSourceURL
	StateTargetURL
	StateEstimation
	StateTableSelect
	StateOptions
	StateMigrationType
	StateMigrating
	StateComplete
	StateError
	StateHistory
)

type Model struct {
	state           State
	sourceURL       string
	targetURL       string
	migrationType   db.MigrationType
	options         db.MigrationOptions
	estimation      *db.EstimationResult
	availableTables []string
	history         []db.MigrationRecord
	textInput       textinput.Model
	progressBar     progress.Model
	spinner         spinner.Model
	progressMsg     string
	currentCommand  string
	progressPct     float64
	progressChan    chan db.ProgressUpdate
	cursor          int
	selectedIndex   int
	scrollOffset    int
	selectedTables  map[string]bool
	loadingTick     int
	errorMsg        string
	successMsg      string
	quitting        bool
	missingDeps     []string
	logo            string
	finalStats      *db.MigrationStats
	systemInfo      *pkgmgr.SystemInfo
}

func InitialModel(logo string) Model {
	ti := textinput.New()
	ti.Placeholder = "postgres://user:password@host:port/dbname"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60

	prog := progress.New(
		progress.WithGradient("#00BFFF", "#FF00FF"),
		progress.WithWidth(60),
		progress.WithoutPercentage(),
	)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	sysInfo := pkgmgr.GetSystemInfo()

	return Model{
		state:          StateCheckingDeps,
		textInput:      ti,
		progressBar:    prog,
		spinner:        s,
		logo:           logo,
		selectedTables: make(map[string]bool),
		systemInfo:     sysInfo,
		options: db.MigrationOptions{
			ParallelJobs: sysInfo.RecommendedWorkers,
			AutoBackup:   true,
		},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.spinner.Tick,
		checkDepsCmd,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

		switch m.state {
		case StateCheckingDeps:
			if len(m.missingDeps) > 0 {
				if msg.String() == "y" || msg.String() == "Y" || msg.String() == "enter" {
					m.state = StateInstallingDeps
					return m, installDepsCmd
				} else if msg.String() == "n" || msg.String() == "N" {
					m.state = StateIntro
				}
			}
		case StateInstallingDeps:
		case StateIntro:
			return m.handleIntro(msg)
		case StateSourceURL:
			return m.handleSourceURL(msg)
		case StateTargetURL:
			return m.handleTargetURL(msg)
		case StateEstimation:
			return m.handleEstimation(msg)
		case StateTableSelect:
			return m.handleTableSelect(msg)
		case StateOptions:
			return m.handleOptions(msg)
		case StateMigrationType:
			return m.handleMigrationType(msg)
		case StateHistory:
			return m.handleHistory(msg)
		case StateComplete:
			if msg.String() == "q" {
				return m, tea.Quit
			}
		case StateError:
			if msg.String() == "q" {
				return m, tea.Quit
			}
		}

	case DepsCheckedMsg:
		if len(msg.Missing) > 0 {
			m.missingDeps = msg.Missing
		} else {
			m.state = StateIntro
		}
		return m, nil

	case DepsInstalledMsg:
		if msg.Err != nil {
			m.state = StateError
			m.errorMsg = "Failed to install dependencies: " + msg.Err.Error()
			return m, tea.Quit
		}
		m.state = StateIntro
		return m, nil

	case EstimationMsg:
		m.estimation = msg.Result
		if msg.Err != nil {
		}
		return m, nil

	case TablesMsg:
		m.availableTables = msg.Tables
		return m, nil

	case HistoryMsg:
		m.history = msg.History
		return m, nil

	case ProgressMsg:
		m.progressPct = msg.Percentage
		m.progressMsg = msg.Message
		if msg.Command != "" {
			m.currentCommand = msg.Command
		}
		if msg.Stats != nil {
			m.finalStats = msg.Stats
		}
		if msg.Percentage < 0 {
			m.state = StateError
			m.errorMsg = msg.Message
			m.progressChan = nil
			return m, nil
		}
		if msg.Percentage >= 1.0 && msg.Stats != nil {
			m.state = StateComplete
			m.successMsg = "Migration completed successfully!"
			m.progressChan = nil
			return m, nil
		}

		var cmds []tea.Cmd
		if m.progressChan != nil {
			cmds = append(cmds, waitForMigrationUpdate(m.progressChan))
		}
		cmd := m.progressBar.SetPercent(m.progressPct)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)

	case MigrationCompleteMsg:
		m.state = StateComplete
		m.successMsg = "Migration completed successfully!"
		m.progressChan = nil
		return m, nil

	case MigrationErrorMsg:
		m.state = StateError
		m.errorMsg = string(msg)
		m.progressChan = nil
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.progressBar.Update(msg)
		m.progressBar = progressModel.(progress.Model)
		return m, cmd

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case TickMsg:
		if m.state == StateMigrating {
			m.loadingTick++
			return m, tickCmd()
		}
		return m, nil
	}

	return m, nil
}
