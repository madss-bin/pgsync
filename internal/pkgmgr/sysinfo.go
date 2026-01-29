package pkgmgr

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type DiskType string

const (
	DiskSSD     DiskType = "SSD"
	DiskHDD     DiskType = "HDD"
	DiskUnknown DiskType = "Unknown"
)

type SystemInfo struct {
	CPUCores           int
	DiskType           DiskType
	RecommendedWorkers int
	Rationale          string
}

func GetCPUCores() int {
	return runtime.NumCPU()
}

func GetDiskType() DiskType {
	entries, err := os.ReadDir("/sys/block")
	if err != nil {
		return DiskUnknown
	}

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, "loop") || strings.HasPrefix(name, "ram") {
			continue
		}

		rotPath := filepath.Join("/sys/block", name, "queue/rotational")
		data, err := os.ReadFile(rotPath)
		if err != nil {
			continue
		}

		rotational := strings.TrimSpace(string(data))
		switch rotational {
		case "0":
			return DiskSSD
		case "1":
			return DiskHDD
		}
	}

	return DiskUnknown
}

func GetRecommendedWorkers() int {
	cores := GetCPUCores()
	diskType := GetDiskType()

	workers := (cores * 3) / 4
	if workers < 2 {
		workers = 2
	}
	if workers > 12 {
		workers = 12
	}

	switch diskType {
	case DiskHDD:
		workers = (workers * 3) / 4
		if workers < 2 {
			workers = 2
		}
	case DiskSSD:
		workers = workers + 1
		if workers > 12 {
			workers = 12
		}
	}

	return workers
}
func GetSystemInfo() *SystemInfo {
	cores := GetCPUCores()
	diskType := GetDiskType()
	workers := GetRecommendedWorkers()

	var rationale string
	switch diskType {
	case DiskSSD:
		rationale = "SSD detected, optimized for parallel I/O"
	case DiskHDD:
		rationale = "HDD detected, reduced parallelism to avoid seek overhead"
	default:
		rationale = "Balanced for general use"
	}

	return &SystemInfo{
		CPUCores:           cores,
		DiskType:           diskType,
		RecommendedWorkers: workers,
		Rationale:          rationale,
	}
}
