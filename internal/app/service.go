package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"yuluwallpaper/internal/config"
	"yuluwallpaper/internal/wallpaper"
)

const WallpaperURL = "http://yulu.frogchou.com/api/v1/quotes/getdesktoppic"

type Service struct {
	mu          sync.Mutex
	cfg         config.Config
	assetsDir   string
	currentPath string

	refreshCh chan struct{}
	updateCh  chan config.Config
	stopCh    chan struct{}
}

func NewService(cfg config.Config, assetsDir string) *Service {
	return &Service{
		cfg:       cfg,
		assetsDir: assetsDir,
		refreshCh: make(chan struct{}, 1),
		updateCh:  make(chan config.Config, 1),
		stopCh:    make(chan struct{}),
	}
}

func (s *Service) Run() {
	s.refresh()

	interval := config.IntervalDuration(s.cfg.IntervalMinutes)
	if interval <= 0 {
		interval = config.IntervalDuration(config.Default().IntervalMinutes)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.refresh()
		case <-s.refreshCh:
			s.refresh()
		case newCfg := <-s.updateCh:
			oldLayout := s.cfg.Layout
			s.cfg = config.Normalize(newCfg)
			ticker.Stop()
			ticker = time.NewTicker(config.IntervalDuration(s.cfg.IntervalMinutes))
			if oldLayout != s.cfg.Layout && s.currentPath != "" {
				if err := wallpaper.Set(s.currentPath, wallpaper.Layout(s.cfg.Layout)); err != nil {
					log.Printf("apply layout failed: %v", err)
				}
			}
		case <-s.stopCh:
			return
		}
	}
}

func (s *Service) RequestRefresh() {
	select {
	case s.refreshCh <- struct{}{}:
	default:
	}
}

func (s *Service) UpdateConfig(cfg config.Config) {
	select {
	case s.updateCh <- cfg:
	default:
		select {
		case <-s.updateCh:
		default:
		}
		s.updateCh <- cfg
	}
}

func (s *Service) Stop() {
	close(s.stopCh)
}

func (s *Service) refresh() {
	path, err := downloadImage(WallpaperURL, s.assetsDir)
	if err != nil {
		log.Printf("download failed: %v", err)
		return
	}
	if err := wallpaper.Set(path, wallpaper.Layout(s.cfg.Layout)); err != nil {
		log.Printf("set wallpaper failed: %v", err)
		return
	}

	s.mu.Lock()
	s.currentPath = path
	s.mu.Unlock()
}

func downloadImage(url, destDir string) (string, error) {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	ext := extensionFromContentType(resp.Header.Get("Content-Type"))
	if ext == "" {
		ext = ".img"
	}

	tmp, err := os.CreateTemp(destDir, "wallpaper-*")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
	}()

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}

	finalPath := filepath.Join(destDir, "wallpaper"+ext)
	_ = os.Remove(finalPath)
	if err := os.Rename(tmp.Name(), finalPath); err != nil {
		return "", err
	}
	return finalPath, nil
}

func extensionFromContentType(contentType string) string {
	contentType = strings.ToLower(contentType)
	switch {
	case strings.Contains(contentType, "image/png"):
		return ".png"
	case strings.Contains(contentType, "image/jpeg"):
		return ".jpg"
	case strings.Contains(contentType, "image/jpg"):
		return ".jpg"
	case strings.Contains(contentType, "image/bmp"):
		return ".bmp"
	case strings.Contains(contentType, "image/gif"):
		return ".gif"
	default:
		return ""
	}
}
