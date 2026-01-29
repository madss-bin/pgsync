package ui

import (
	"fmt"
	"math"
	"strings"

	"pgsync/internal/db"

	"github.com/charmbracelet/lipgloss"
)

var funMessages = []string{
	"we are cookin'",
	"hold tight...",
	"migrating bits...",
	"we balls fr",
	"postgres go brrr",
	"optimizing reality...",
	"asking postgres nicely",
	"casting pg_dump...",
	"trust the process",
	"almost there...",
	"definitely not stuck",
	"pure elegance...",
}

func (m Model) viewProgress() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(ProgressTextStyle.Render(m.progressMsg))
	b.WriteString("\n")

	if m.currentCommand != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("> " + m.currentCommand))
		b.WriteString("\n")
	} else {
		b.WriteString("\n")
	}
	b.WriteString("\n")

	b.WriteString(m.progressBar.ViewAs(m.progressPct))
	b.WriteString(" ")

	percent := math.Round(m.progressPct * 100)
	pctStr := fmt.Sprintf("%.0f%%", percent)

	msgIndex := (m.loadingTick / 40) % len(funMessages)
	funMsg := funMessages[msgIndex]
	if percent >= 100 {
		funMsg = "Done!"
	}

	b.WriteString(HintStyle.Render(fmt.Sprintf("%s %s", pctStr, funMsg)))
	b.WriteString("\n\n")

	return b.String()
}

func (m Model) viewComplete() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(colorNeonBlue).Render(m.logo))
	b.WriteString("\n\n")
	b.WriteString(SuccessStyle.Render("✓ " + m.successMsg))
	b.WriteString("\n\n")

	b.WriteString(PromptStyle.Render("Migration Summary"))
	b.WriteString("\n\n")

	if m.finalStats != nil {
		modeStr := "Full (Schema + Data)"
		switch m.finalStats.MigrationType {
		case db.SchemaOnly:
			modeStr = "Schema Only"
		case db.DataOnly:
			modeStr = "Data Only"
		}
		b.WriteString(fmt.Sprintf("   Mode:           %s\n", modeStr))

		if m.finalStats.TablesMigrated > 0 {
			b.WriteString(fmt.Sprintf("   Tables:         %d migrated\n", m.finalStats.TablesMigrated))
		} else {
			b.WriteString("   Tables:         All tables\n")
		}

		b.WriteString(fmt.Sprintf("   Duration:       %s\n", m.finalStats.Duration))

		if m.finalStats.BackupPath != "" {
			b.WriteString(fmt.Sprintf("   Backup:         %s\n", m.finalStats.BackupPath))
		}

		if len(m.finalStats.Warnings) > 0 {
			b.WriteString("\n")
			b.WriteString(WarningStyle.Render("   ⚠ Warnings:"))
			b.WriteString("\n")
			for _, w := range m.finalStats.Warnings {
				b.WriteString(fmt.Sprintf("     • %s\n", w))
			}
		}

		if m.finalStats.DidRollback {
			b.WriteString("\n")
			if m.finalStats.RollbackSuccess {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render("   ⟲ Rollback was performed (target restored)"))
			} else {
				b.WriteString(ErrorStyle.Render("   ⟲ Rollback attempted but failed"))
			}
			b.WriteString("\n")
		}

		if m.finalStats.LogPath != "" {
			b.WriteString("\n")
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(fmt.Sprintf("   Log saved to: %s", m.finalStats.LogPath)))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("press q to exit"))
	b.WriteString("\n\n")

	return b.String()
}

func (m Model) viewError() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(ErrorStyle.Render("✗ Migration failed"))
	b.WriteString("\n\n")
	b.WriteString(ErrorMessageStyle.Render(m.errorMsg))
	b.WriteString("\n\n")

	if m.finalStats != nil {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Details:"))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("   Duration: %s\n", m.finalStats.Duration))

		if m.finalStats.BackupPath != "" {
			b.WriteString(fmt.Sprintf("   Backup:   %s\n", m.finalStats.BackupPath))
		}

		if m.finalStats.DidRollback {
			b.WriteString("\n")
			if m.finalStats.RollbackSuccess {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render("   ⟲ Rollback successful - target database restored"))
			} else {
				b.WriteString(WarningStyle.Render("   ⟲ Rollback also failed - manual intervention may be needed"))
			}
			b.WriteString("\n")
		}

		if m.finalStats.LogPath != "" {
			b.WriteString("\n")
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(fmt.Sprintf("   Log saved to: %s", m.finalStats.LogPath)))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("press q to exit"))
	b.WriteString("\n\n")
	return b.String()
}
