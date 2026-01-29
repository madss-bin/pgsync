package ui

import (
	"fmt"
	"strings"

	"pgsync/internal/db"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) viewCheckingDeps() string {
	var b strings.Builder
	b.WriteString("\n")

	if len(m.missingDeps) > 0 {
		b.WriteString(WarningStyle.Render(fmt.Sprintf("PostgreSQL tools not found: %s", strings.Join(m.missingDeps, ", "))))
		b.WriteString("\n\n")
		b.WriteString(PromptStyle.Render("Install now? [Y/n]: "))
	} else {
		b.WriteString("   " + m.spinner.View() + " Checking dependencies...\n")
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("   > which pg_dump pg_restore psql"))
		b.WriteString("\n")
		b.WriteString(m.progressBar.ViewAs(0.1))
		b.WriteString("\n")
	}
	return b.String()
}

func (m Model) viewInstallingDeps() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString("   " + m.spinner.View() + " Installing PostgreSQL client tools...\n")
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("   > package manager install postgresql"))
	b.WriteString("\n")
	b.WriteString(m.progressBar.ViewAs(0.4))
	b.WriteString("\n")
	return b.String()
}

func (m Model) viewIntro() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(colorNeonBlue).Render(m.logo))
	b.WriteString("\n")
	b.WriteString(SubtitleStyle.Render("   PostgreSQL database migration tool"))
	b.WriteString("\n\n")
	b.WriteString(HelpStyle.Render("   press enter to start • h for history"))
	b.WriteString("\n\n")

	return b.String()
}

func (m Model) viewHistory() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(PromptStyle.Render("Migration History"))
	b.WriteString("\n\n")

	if len(m.history) == 0 {
		b.WriteString("   No history found.\n")
	} else {
		for i, h := range m.history {
			if i >= 10 {
				break
			}
			statusColor := "2"
			if h.Status == db.StatusFailed {
				statusColor = "1"
			}

			b.WriteString(fmt.Sprintf("   %s  %s -> %s  (%s)\n",
				lipgloss.NewStyle().Foreground(lipgloss.Color(statusColor)).Render(string(h.Status)),
				h.Timestamp.Format("2006-01-02 15:04"),
				h.Target,
				h.Duration,
			))
		}
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("esc to back"))
	b.WriteString("\n\n")
	return b.String()
}
