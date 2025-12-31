package logger

import (
	"log"
	"os"
	"path/filepath"
)

var logFile *os.File

func Init(path string) error {
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	logFile = file
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	return nil
}

func Close() {
	if logFile != nil {
		_ = logFile.Close()
	}
}