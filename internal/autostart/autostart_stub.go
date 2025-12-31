//go:build !windows && !darwin

package autostart

import "errors"

func enable(appName, execPath, logPath string) error {
	return errors.New("autostart not supported on this platform")
}

func disable(appName string) error {
	return errors.New("autostart not supported on this platform")
}