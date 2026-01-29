package db

import (
	"fmt"
	"strings"
)

func ValidateURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	if !strings.HasPrefix(url, "postgres://") && !strings.HasPrefix(url, "postgresql://") {
		return fmt.Errorf("please enter a valid PostgreSQL URL (must start with postgres:// or postgresql://)")
	}

	return nil
}

func URLsAreDifferent(source, target string) error {
	if source == target {
		return fmt.Errorf("source and target URLs cannot be the same")
	}
	return nil
}
