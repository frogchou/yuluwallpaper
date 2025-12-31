package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

const AppName = "yuluwallpaper"

type Layout string

const (
	LayoutTile    Layout = "tile"
	LayoutStretch Layout = "stretch"
	LayoutFit     Layout = "fit"
	LayoutFill    Layout = "fill"
	LayoutCenter  Layout = "center"
)

type Config struct {
	IntervalMinutes int    `json:"interval_minutes"`
	Layout          Layout `json:"layout"`
	AutoStart       bool   `json:"auto_start"`
}

type IntervalOption struct {
	Label   string
	Minutes int
}

var intervalOptions = []IntervalOption{
	{Label: "10分钟", Minutes: 10},
	{Label: "20分钟", Minutes: 20},
	{Label: "30分钟", Minutes: 30},
	{Label: "40分钟", Minutes: 40},
	{Label: "50分钟", Minutes: 50},
	{Label: "1小时", Minutes: 60},
	{Label: "2小时", Minutes: 120},
	{Label: "4小时", Minutes: 240},
	{Label: "8小时", Minutes: 480},
	{Label: "1天", Minutes: 1440},
	{Label: "2天", Minutes: 2880},
}

func IntervalOptions() []IntervalOption {
	return intervalOptions
}

func IntervalLabel(minutes int) string {
	for _, opt := range intervalOptions {
		if opt.Minutes == minutes {
			return opt.Label
		}
	}
	return ""
}

func IntervalDuration(minutes int) time.Duration {
	return time.Duration(minutes) * time.Minute
}

func Default() Config {
	return Config{
		IntervalMinutes: 60,
		Layout:          LayoutFill,
		AutoStart:       false,
	}
}

func Normalize(cfg Config) Config {
	if !validInterval(cfg.IntervalMinutes) {
		cfg.IntervalMinutes = Default().IntervalMinutes
	}
	switch cfg.Layout {
	case LayoutTile, LayoutStretch, LayoutFit, LayoutFill, LayoutCenter:
	default:
		cfg.Layout = Default().Layout
	}
	return cfg
}

func AppDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	if base == "" {
		return "", errors.New("user config dir is empty")
	}
	return filepath.Join(base, AppName), nil
}

func ConfigPath() (string, error) {
	dir, err := AppDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func AssetsDir() (string, error) {
	dir, err := AppDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "assets"), nil
}

func LogPath() (string, error) {
	dir, err := AppDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "yuluwallpaper.log"), nil
}

func Load() (Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return Default(), err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return Default(), err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Default(), err
	}
	return Normalize(cfg), nil
}

func Save(cfg Config) error {
	cfg = Normalize(cfg)
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func validInterval(minutes int) bool {
	for _, opt := range intervalOptions {
		if opt.Minutes == minutes {
			return true
		}
	}
	return false
}
