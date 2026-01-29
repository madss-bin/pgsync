<h1 align="center">pgsync</h1>

<div align="center">

[![Go](https://img.shields.io/badge/Go-1.22+-007D9C?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-TUI-000000?style=for-the-badge&logo=bubbletea)](https://github.com/charmbracelet/bubbletea)
[![Lip Gloss](https://img.shields.io/badge/Lip%20Gloss-Styling-FF00FF?style=for-the-badge&logo=lipgloss)](https://github.com/charmbracelet/lipgloss)
[![Cobra](https://img.shields.io/badge/Cobra-CLI-FF0000?style=for-the-badge&logo=cobra)](https://github.com/spf13/cobra)
<br>
<br>

</div>

## Features

- **Pre-Flight Checks**: Validates versions, extensions, and disk space before migration begins.
- **Table Selection**: Interactive UI to include or exclude specific tables.
- **Safety Backups**: Optional auto-backup of target database before overwriting.
- **Rollback on Failure**: Automatically restores from backup if migration fails.
- **Smart Parallelism**: Detects CPU cores and disk type to recommend optimal worker count.
- **Migration Summary**: Displays a complete recap after migration with mode, duration, and warnings.
- **History Tracking**: Logs previous migrations with status and duration.
- **Auto-Detection**: Installs required PostgreSQL client tools if missing.

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/madss-bin/pgsync/main/install.sh | bash
```

- **Arch Linux**: Uses `makepkg` to build and install a system package.
- **Other Distros**: Builds from source and installs to `/usr/local/bin`.

## Usage

Run the command to launch the wizard:

```bash
pgsync
```

Follow the on-screen prompts to:

1. Enter source and target database URLs
2. Review pre-flight checks (version compatibility, extensions, size)
3. Select tables to migrate (or migrate all)
4. Configure parallel workers and safety backup
5. Choose migration type (full, schema-only, or data-only)

After completion, the summary screen shows:

- Migration mode and duration
- Tables migrated
- Backup file location
- Any warnings encountered

## Requirements

- **Linux** (Arch, Fedora, Ubuntu/Debian supported for auto-setup)
- **PostgreSQL** (Client tools required, auto-installed if missing)
- **Go 1.22+** (Required for building from source)

---

## Manual Commands

For those who prefer the terminal directly:

**Full Migration**
```bash
# Dump source
pg_dump "postgres://user:pass@source_host/db" -Fc -f dump.pkg --no-owner --no-privileges --verbose

# Restore to target with parallel workers
pg_restore -d "postgres://user:pass@target_host/db" -j 4 -c --if-exists --no-owner --no-privileges --verbose dump.pkg
```

**Schema Only**
Add `--schema-only` to the `pg_dump` command.

**Data Only**
Add `--data-only` to the `pg_dump` command.

---

## License

[MIT](https://github.com/madss-bin/pgsync/raw/main/LICENSE)

## Troubleshooting

**Migration Logs**  
Every migration creates a detailed log file (path displayed on completion). If a migration fails or leaves the database empty, check this log for `pg_dump` or `pg_restore` errors.

**Supabase Users**  
Do not use the Transaction Pooler (port 6543). You must use the Session Mode connection (port 5432) for migrations to work correctly. The tool will block port 6543 to prevent issues.

### Uninstall
To remove pgsync:
```bash
./uninstall.sh
```