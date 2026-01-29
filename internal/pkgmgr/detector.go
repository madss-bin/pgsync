package pkgmgr

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

type Distro string

const (
	Arch    Distro = "arch"
	Fedora  Distro = "fedora"
	Ubuntu  Distro = "ubuntu"
	Unknown Distro = "unknown"
)

func DetectDistro() Distro {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return Unknown
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var id string
	var idLike string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ID=") {
			id = strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
		}
		if strings.HasPrefix(line, "ID_LIKE=") {
			idLike = strings.Trim(strings.TrimPrefix(line, "ID_LIKE="), "\"")
		}
	}

	switch strings.ToLower(id) {
	case "arch", "archlinux", "manjaro":
		return Arch
	case "fedora", "rhel", "centos":
		return Fedora
	case "ubuntu", "debian":
		return Ubuntu
	}

	idLikeLower := strings.ToLower(idLike)
	if strings.Contains(idLikeLower, "arch") {
		return Arch
	}
	if strings.Contains(idLikeLower, "fedora") || strings.Contains(idLikeLower, "rhel") {
		return Fedora
	}
	if strings.Contains(idLikeLower, "debian") || strings.Contains(idLikeLower, "ubuntu") {
		return Ubuntu
	}

	return Unknown
}

func CheckDependencies() []string {
	missing := []string{}
	tools := []string{"pg_dump", "pg_restore", "psql"}

	for _, tool := range tools {
		if _, err := exec.LookPath(tool); err != nil {
			missing = append(missing, tool)
		}
	}
	return missing
}
