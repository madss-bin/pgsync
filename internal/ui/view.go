package ui

func (m Model) View() string {
	if m.quitting {
		return "\n" + WarningStyle.Render("Exiting...") + "\n\n"
	}

	switch m.state {
	case StateCheckingDeps:
		return m.viewCheckingDeps()
	case StateInstallingDeps:
		return m.viewInstallingDeps()
	case StateIntro:
		return m.viewIntro()
	case StateHistory:
		return m.viewHistory()
	case StateSourceURL:
		return m.viewSourceURL()
	case StateTargetURL:
		return m.viewTargetURL()
	case StateEstimation:
		return m.viewEstimation()
	case StateTableSelect:
		return m.viewTableSelect()
	case StateOptions:
		return m.viewOptions()
	case StateMigrationType:
		return m.viewMigrationType()
	case StateMigrating:
		return m.viewProgress()
	case StateComplete:
		return m.viewComplete()
	case StateError:
		return m.viewError()
	}

	return ""
}
