package db

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetTables(url string) ([]string, error) {
	query := "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_name;"

	cmd := exec.Command("psql", url, "-t", "-c", query)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %v (output: %s)", err, string(out))
	}

	var tables []string
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed != "" {
			tables = append(tables, trimmed)
		}
	}

	return tables, nil
}
