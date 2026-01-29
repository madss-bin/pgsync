package ui

import (
	"fmt"
	"strings"

	"pgsync/internal/db"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) viewEstimation() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(PromptStyle.Render("Pre-flight Checks"))
	b.WriteString("\n\n")

	if m.estimation == nil {
		b.WriteString("   " + m.spinner.View() + " Analyzing databases...\n")
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("   > psql ... SHOW server_version"))
		b.WriteString("\n")
		b.WriteString(m.progressBar.ViewAs(0.2))
		b.WriteString("\n")
		return b.String()
	}

	b.WriteString(fmt.Sprintf("   Source: %s\n", m.estimation.SourceVersion))
	b.WriteString(fmt.Sprintf("   Target: %s\n\n", m.estimation.TargetVersion))

	for _, c := range m.estimation.Checks {
		icon := "✓"
		color := "2"
		switch c.Status {
		case db.StatusYellow:
			icon = "⚠"
			color = "3"
		case db.StatusRed:
			icon = "✗"
			color = "1"
		}

		line := fmt.Sprintf("%s %s: %s", icon, c.Name, c.Message)
		b.WriteString("   " + lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(line) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("enter to continue"))
	b.WriteString("\n\n")
	return b.String()
}

func (m Model) viewTableSelect() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(PromptStyle.Render("Select Tables (Space to toggle, A for all)"))
	b.WriteString("\n\n")

	if m.availableTables == nil {
		b.WriteString("   " + m.spinner.View() + " Fetching tables...\n")
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("   > psql ... FROM information_schema.tables"))
		b.WriteString("\n")
		b.WriteString(m.progressBar.ViewAs(0.3))
		b.WriteString("\n")
		return b.String()
	}

	start := m.scrollOffset
	end := start + 10
	if end > len(m.availableTables) {
		end = len(m.availableTables)
	}

	for i := start; i < end; i++ {
		table := m.availableTables[i]
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := "[ ]"
		if m.selectedTables[table] {
			checked = "[x]"
		}

		style := UnselectedItemStyle
		if m.cursor == i {
			style = SelectedItemStyle
		}

		b.WriteString(style.Render(fmt.Sprintf("%s %s %s", cursor, checked, table)) + "\n")
	}

	if len(m.availableTables) > 10 {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(fmt.Sprintf("\n   ... %d more tables (↓ to scroll)", len(m.availableTables)-end)))
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("enter to confirm selection"))
	b.WriteString("\n\n")
	return b.String()
}

func (m Model) viewOptions() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(PromptStyle.Render("Migration Options"))
	b.WriteString("\n\n")
	if m.systemInfo != nil {
		infoLine := fmt.Sprintf("   ℹ Detected: %d CPU cores, %s", m.systemInfo.CPUCores, m.systemInfo.DiskType)
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(infoLine))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(fmt.Sprintf("   %s", m.systemInfo.Rationale)))
		b.WriteString("\n\n")
	}

	style := UnselectedItemStyle
	cursor := " "
	if m.cursor == 0 {
		style = SelectedItemStyle
		cursor = ">"
	}

	jobsInfo := fmt.Sprintf("Parallel Jobs: %d", m.options.ParallelJobs)
	if m.systemInfo != nil && m.options.ParallelJobs == m.systemInfo.RecommendedWorkers {
		jobsInfo += " (recommended)"
	}
	if m.cursor == 0 {
		jobsInfo += "  (←/→ to change)"
	}
	b.WriteString(style.Render(fmt.Sprintf("%s %s", cursor, jobsInfo)))
	b.WriteString("\n")

	style = UnselectedItemStyle
	cursor = " "
	if m.cursor == 1 {
		style = SelectedItemStyle
		cursor = ">"
	}

	backupStr := "No"
	if m.options.AutoBackup {
		backupStr = "Yes"
	}
	backupInfo := fmt.Sprintf("Safety Backup: %s", backupStr)
	if m.cursor == 1 {
		backupInfo += " (Space to toggle)"
	}
	b.WriteString(style.Render(fmt.Sprintf("%s %s", cursor, backupInfo)))
	b.WriteString("\n\n")

	style = UnselectedItemStyle
	cursor = " "
	if m.cursor == 2 {
		style = SelectedItemStyle
		cursor = ">"
	}
	b.WriteString(style.Render(fmt.Sprintf("%s Continue to Confirmation", cursor)))
	b.WriteString("\n\n")

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("↑/↓ or j/k to move • enter to confirm"))
	b.WriteString("\n\n")
	return b.String()
}

func (m Model) viewMigrationType() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(PromptStyle.Render("Select migration type:"))
	b.WriteString("\n\n")

	choices := []string{
		"Schema + Data (full migration)",
		"Schema only (no data)",
		"Data only (no schema)",
	}

	for i, choice := range choices {
		style := UnselectedItemStyle
		cursor := " "
		if m.selectedIndex == i {
			style = SelectedItemStyle
			cursor = ">"
		}
		b.WriteString(style.Render(fmt.Sprintf("%s %s", cursor, choice)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("↑/↓ or j/k to move • enter to confirm"))
	b.WriteString("\n\n")
	return b.String()
}
