//go:build windows

package autostart

import (
	"strings"

	"golang.org/x/sys/windows/registry"
)

func enable(appName, execPath, logPath string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	value := execPath
	if strings.ContainsAny(execPath, " ") {
		value = `"` + execPath + `"`
	}
	return key.SetStringValue(appName, value)
}

func disable(appName string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	if err := key.DeleteValue(appName); err != nil && err != registry.ErrNotExist {
		return err
	}
	return nil
}