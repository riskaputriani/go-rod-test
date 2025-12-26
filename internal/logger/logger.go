package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Logger adalah struct untuk mengelola logging
type Logger struct {
	writer io.Writer
	path   string
}

// New membuat instance logger baru
func New() (*Logger, error) {
	logFile, logPath, err := openLogFile()
	if err != nil {
		// Jika gagal membuka file, gunakan stdout saja
		return &Logger{
			writer: os.Stdout,
			path:   "",
		}, fmt.Errorf("failed to open log file: %w", err)
	}

	return &Logger{
		writer: io.MultiWriter(os.Stdout, logFile),
		path:   logPath,
	}, nil
}

// LogKV mencatat key-value pair
func (l *Logger) LogKV(key, value string) {
	fmt.Fprintf(l.writer, "%s: %s\n", key, value)
}

// GetPath mengembalikan path file log
func (l *Logger) GetPath() string {
	return l.path
}

// GetWriter mengembalikan writer
func (l *Logger) GetWriter() io.Writer {
	return l.writer
}

// SanitizeValue membersihkan nilai dari karakter newline
func SanitizeValue(value string) string {
	value = strings.ReplaceAll(value, "\r", "\\r")
	value = strings.ReplaceAll(value, "\n", "\\n")
	return value
}

// openLogFile membuka atau membuat file log
func openLogFile() (*os.File, string, error) {
	// Cek environment variable
	if path := strings.TrimSpace(os.Getenv("RUNTIME_LOG_PATH")); path != "" {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, "", err
		}
		file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		return file, path, err
	}

	// Coba di working directory
	if cwd, err := os.Getwd(); err == nil {
		path := filepath.Join(cwd, "runtime-info.log")
		file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err == nil {
			return file, path, nil
		}
	}

	// Fallback ke temp directory
	path := filepath.Join(os.TempDir(), "runtime-info.log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	return file, path, err
}
