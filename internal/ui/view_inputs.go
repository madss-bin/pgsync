package ui

import (
	"strings"
)

func (m Model) viewSourceURL() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(PromptStyle.Render("Enter source database URL:"))
	b.WriteString("\n\n")
	b.WriteString(m.textInput.View())
	b.WriteString("\n")
	if m.errorMsg != "" {
		b.WriteString("\n")
		b.WriteString(ErrorMessageStyle.Render("✗ " + m.errorMsg))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("enter to continue"))
	b.WriteString("\n\n")
	return b.String()
}

func (m Model) viewTargetURL() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(PromptStyle.Render("Enter target database URL:"))
	b.WriteString("\n\n")
	b.WriteString(m.textInput.View())
	b.WriteString("\n")
	if m.errorMsg != "" {
		b.WriteString("\n")
		b.WriteString(ErrorMessageStyle.Render("✗ " + m.errorMsg))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("enter to continue"))
	b.WriteString("\n\n")
	return b.String()
}
