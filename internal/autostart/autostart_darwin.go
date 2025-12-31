//go:build darwin

package autostart

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func enable(appName, execPath, logPath string) error {
	label := labelFor(appName)
	plistPath, err := launchAgentPath(label)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(plistPath), 0o755); err != nil {
		return err
	}

	content := fmt.Sprintf(plistTemplate,
		label,
		xmlEscape(execPath),
		xmlEscape(logPath),
		xmlEscape(logPath),
	)

	if err := os.WriteFile(plistPath, []byte(content), 0o644); err != nil {
		return err
	}

	_ = exec.Command("launchctl", "load", "-w", plistPath).Run()
	return nil
}

func disable(appName string) error {
	label := labelFor(appName)
	plistPath, err := launchAgentPath(label)
	if err != nil {
		return err
	}

	_ = exec.Command("launchctl", "unload", "-w", plistPath).Run()
	if err := os.Remove(plistPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func launchAgentPath(label string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Library", "LaunchAgents", label+".plist"), nil
}

func labelFor(appName string) string {
	clean := strings.ToLower(appName)
	clean = strings.ReplaceAll(clean, " ", "")
	clean = strings.ReplaceAll(clean, ".", "")
	if clean == "" {
		clean = "yuluwallpaper"
	}
	return "com.yulu." + clean
}

func xmlEscape(value string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(value)
}

const plistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>%s</string>
  <key>ProgramArguments</key>
  <array>
    <string>%s</string>
  </array>
  <key>RunAtLoad</key>
  <true/>
  <key>StandardOutPath</key>
  <string>%s</string>
  <key>StandardErrorPath</key>
  <string>%s</string>
</dict>
</plist>
`