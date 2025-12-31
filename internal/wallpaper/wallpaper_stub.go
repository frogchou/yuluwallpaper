//go:build !windows && !darwin

package wallpaper

import "errors"

func setWallpaper(path string, layout Layout) error {
	return errors.New("wallpaper not supported on this platform")
}