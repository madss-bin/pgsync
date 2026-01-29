package cmd

import (
	"fmt"
	"os"
	"pgsync/internal/db"
	"pgsync/internal/pkgmgr"
	"pgsync/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pgsync",
	Short: "PostgreSQL database migration tool",
	Long:  `Launch an interactive session to migrate a PostgreSQL database from source to target.`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Version: "1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		if err := db.CheckDependencies(); err != nil {
			fmt.Println(ui.WarningStyle.Render("\n⚠ PostgreSQL tools not found"))
			fmt.Print("\nInstall now? [Y/n]: ")

			var response string
			fmt.Scanln(&response)

			if response == "n" || response == "N" {
				fmt.Println(ui.ErrorStyle.Render("\nCancelled. Please install PostgreSQL client tools manually."))
				os.Exit(1)
			}

			fmt.Println(ui.PromptStyle.Render("\n→ Installing PostgreSQL client tools..."))
			if err := pkgmgr.InstallPostgreSQL(); err != nil {
				fmt.Println(ui.ErrorStyle.Render("✗ Installation failed: " + err.Error()))
				os.Exit(1)
			}
			fmt.Println(ui.SuccessStyle.Render("✓ Installation complete\n"))
		}
		p := tea.NewProgram(ui.InitialModel(logo), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Println(ui.ErrorStyle.Render("Error: " + err.Error()))
			os.Exit(1)
		}
	},
}

var logo string

func Execute(l string) {
	logo = l
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
