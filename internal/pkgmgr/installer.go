package pkgmgr

import (
	"fmt"
	"os/exec"
)

func InstallPostgreSQL() error {
	distro := DetectDistro()
	
	var cmd *exec.Cmd
	
	switch distro {
	case Arch:
		cmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", "postgresql")
	case Fedora:
		cmd = exec.Command("sudo", "dnf", "install", "-y", "postgresql")
	case Ubuntu:
		updateCmd := exec.Command("sudo", "apt", "update")
		if err := updateCmd.Run(); err != nil {
			return fmt.Errorf("failed to update package list: %w", err)
		}
		cmd = exec.Command("sudo", "apt", "install", "-y", "postgresql-client")
	default:
		return fmt.Errorf("unsupported distribution - please install PostgreSQL client tools manually")
	}
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}
	
	return nil
}

func GetInstallCommand(distro Distro) string {
	switch distro {
	case Arch:
		return "sudo pacman -S postgresql"
	case Fedora:
		return "sudo dnf install postgresql"
	case Ubuntu:
		return "sudo apt update && sudo apt install postgresql-client"
	default:
		return "Please install PostgreSQL client tools (pg_dump and psql)"
	}
}

