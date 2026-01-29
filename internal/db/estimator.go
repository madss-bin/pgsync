package db

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type CheckStatus string

const (
	StatusGreen  CheckStatus = "green"
	StatusYellow CheckStatus = "yellow"
	StatusRed    CheckStatus = "red"
)

type CheckResult struct {
	Name    string
	Status  CheckStatus
	Message string
}

type EstimationResult struct {
	SourceVersion string
	TargetVersion string
	DbSize        string
	TableCount    int
	Checks        []CheckResult
}

const cmdTimeout = 10 * time.Second

func Estimate(source, target string) (*EstimationResult, error) {
	res := &EstimationResult{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	addError := func(err error) {
		mu.Lock()
		errs = append(errs, err)
		mu.Unlock()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := CheckConnection(source); err != nil {
			addError(fmt.Errorf("source unreachable: %w", err))
			return
		}

		ver, err := getPGVersion(source)
		if err != nil {
			addError(fmt.Errorf("failed to get source version: %w", err))
			return
		}
		mu.Lock()
		res.SourceVersion = ver
		mu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := CheckConnection(target); err != nil {
			addError(fmt.Errorf("target unreachable: %w", err))
			return
		}

		ver, err := getPGVersion(target)
		if err != nil {
			addError(fmt.Errorf("failed to get target version: %w", err))
			return
		}
		mu.Lock()
		res.TargetVersion = ver
		mu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		size, err := getDBSize(source)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			res.Checks = append(res.Checks, CheckResult{Name: "Disk Space", Status: StatusYellow, Message: "Could not estimate size"})
		} else {
			res.DbSize = size
			res.Checks = append(res.Checks, CheckResult{Name: "Estimated Size", Status: StatusGreen, Message: size})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err := getTableCount(source)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			res.Checks = append(res.Checks, CheckResult{Name: "Table Count", Status: StatusYellow, Message: "Unknown"})
		} else {
			res.TableCount = count
			res.Checks = append(res.Checks, CheckResult{Name: "Objects", Status: StatusGreen, Message: fmt.Sprintf("%d tables", count)})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		exts, err := getExtensions(source)
		var tgtExts []string
		if err == nil {
			tgtExts, err = getExtensions(target)
		}

		mu.Lock()
		defer mu.Unlock()

		if err == nil {
			missing := []string{}
			for _, srcExt := range exts {
				found := false
				for _, tgtExt := range tgtExts {
					if srcExt == tgtExt {
						found = true
						break
					}
				}
				if !found {
					missing = append(missing, srcExt)
				}
			}
			if len(missing) > 0 {
				res.Checks = append(res.Checks, CheckResult{Name: "Extensions", Status: StatusRed, Message: fmt.Sprintf("Missing on target: %s", strings.Join(missing, ", "))})
			} else {
				res.Checks = append(res.Checks, CheckResult{Name: "Extensions", Status: StatusGreen, Message: "All extensions present"})
			}
		} else {
			res.Checks = append(res.Checks, CheckResult{Name: "Extensions", Status: StatusYellow, Message: "Could not check extensions"})
		}
	}()

	wg.Wait()

	if len(errs) > 0 {
		return nil, errs[0]
	}

	verCheck := CheckResult{Name: "Postgres Compatibility", Status: StatusGreen, Message: "Versions are compatible"}
	if res.SourceVersion != res.TargetVersion {
		verCheck.Status = StatusYellow
		verCheck.Message = fmt.Sprintf("Version mismatch: %s -> %s", res.SourceVersion, res.TargetVersion)
	}
	res.Checks = append(res.Checks, verCheck)

	if strings.Contains(source, ":6543") || strings.Contains(target, ":6543") {
		res.Checks = append(res.Checks, CheckResult{
			Name:    "Connection Mode",
			Status:  StatusYellow,
			Message: "Port 6543 detected (Supabase Pooler). Use 5432 if possible.",
		})
	}

	return res, nil
}

func getPGVersion(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "psql", url, "-w", "-t", "-c", "SHOW server_version;")
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("connection timed out")
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func getDBSize(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "psql", url, "-w", "-t", "-c", "SELECT pg_size_pretty(pg_database_size(current_database()));")
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("connection timed out")
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func getTableCount(url string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "psql", url, "-w", "-t", "-c", "SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public';")
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return 0, fmt.Errorf("connection timed out")
	}
	if err != nil {
		return 0, err
	}
	var count int
	fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &count)
	return count, nil
}

func getExtensions(url string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "psql", url, "-w", "-t", "-c", "SELECT extname FROM pg_extension;")
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return nil, fmt.Errorf("connection timed out")
	}
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var exts []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed != "" {
			exts = append(exts, trimmed)
		}
	}
	return exts, nil
}
