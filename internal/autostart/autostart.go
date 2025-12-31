package autostart

func Enable(appName, execPath, logPath string) error {
	return enable(appName, execPath, logPath)
}

func Disable(appName string) error {
	return disable(appName)
}