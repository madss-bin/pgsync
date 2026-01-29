package db

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func CheckDependencies() error {
	if _, err := exec.LookPath("pg_dump"); err != nil {
		return fmt.Errorf("pg_dump not found in PATH")
	}

	if _, err := exec.LookPath("psql"); err != nil {
		return fmt.Errorf("psql not found in PATH")
	}

	return nil
}

func CheckConnection(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "pg_isready", "-d", url, "-t", "3")
	output, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("connection check timed out (5s)")
	}

	if err != nil {
		outputStr := strings.TrimSpace(string(output))
		return fmt.Errorf("connection check failed: %s (%w)", outputStr, err)
	}

	return nil
}
