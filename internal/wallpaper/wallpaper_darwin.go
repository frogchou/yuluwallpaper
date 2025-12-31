//go:build darwin

package wallpaper

import (
	"fmt"
	"os/exec"
	"strings"
)

func setWallpaper(path string, layout Layout) error {
	scaling := "stretch to fill"
	switch layout {
	case LayoutTile:
		scaling = "tile"
	case LayoutCenter:
		scaling = "center"
	case LayoutFit:
		scaling = "fit to screen"
	case LayoutFill:
		scaling = "fill screen"
	case LayoutStretch:
		scaling = "stretch to fill"
	}

	script := fmt.Sprintf(`tell application "System Events"
repeat with d in desktops
set picture of d to "%s"
set picture scaling of d to %s
end repeat
end tell`, escapeAppleScriptString(path), scaling)

	cmd := exec.Command("osascript", "-e", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("osascript failed: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func escapeAppleScriptString(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	return value
}
