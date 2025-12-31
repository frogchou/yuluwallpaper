//go:build windows

package wallpaper

import (
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"unsafe"
)

func setWallpaper(path string, layout Layout) error {
	if err := setWindowsStyle(layout); err != nil {
		return err
	}
	ptr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	return systemParametersInfo(spiSetDesktopWallpaper, 0, uintptr(unsafe.Pointer(ptr)), spifUpdateIniFile|spifSendChange)
}

func setWindowsStyle(layout Layout) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\Desktop`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	tile := "0"
	style := "2"
	switch layout {
	case LayoutTile:
		tile = "1"
		style = "0"
	case LayoutCenter:
		tile = "0"
		style = "0"
	case LayoutFit:
		tile = "0"
		style = "6"
	case LayoutFill:
		tile = "0"
		style = "10"
	case LayoutStretch:
		tile = "0"
		style = "2"
	}

	if err := key.SetStringValue("WallpaperStyle", style); err != nil {
		return err
	}
	return key.SetStringValue("TileWallpaper", tile)
}

const (
	spiSetDesktopWallpaper = 0x0014
	spifUpdateIniFile      = 0x0001
	spifSendChange         = 0x0002
)

func systemParametersInfo(action, param uint32, pvParam uintptr, winIni uint32) error {
	user32 := windows.NewLazySystemDLL("user32.dll")
	proc := user32.NewProc("SystemParametersInfoW")
	ret, _, err := proc.Call(uintptr(action), uintptr(param), pvParam, uintptr(winIni))
	if ret == 0 {
		if err == nil {
			err = windows.GetLastError()
		}
		return err
	}
	return nil
}
