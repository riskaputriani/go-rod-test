package browser

import (
	"os"
	"path/filepath"
)

// Config berisi konfigurasi untuk Chrome
type Config struct {
	// URL untuk download Chrome for Testing
	DownloadURL string

	// Nama direktori instalasi
	InstallDirName string

	// Versi Chrome
	Version string

	// URL untuk download dependencies (opsional)
	DependenciesURLs []string
}

// DefaultConfig mengembalikan konfigurasi default untuk Chrome for Testing
func DefaultConfig() Config {
	return Config{
		// Chrome for Testing - sudah include banyak dependencies
		DownloadURL:      "https://storage.googleapis.com/chrome-for-testing-public/131.0.6778.204/linux64/chrome-linux64.zip",
		InstallDirName:   "chrome-for-testing",
		Version:          "131.0.6778.204",
		DependenciesURLs: []string{
			// Dependencies akan didownload jika diperlukan
			// Untuk sekarang kosong, akan ditambahkan jika ada error
		},
	}
}

// NewChromiumManagerWithConfig membuat instance baru dengan config custom
func NewChromiumManagerWithConfig(config Config, logger func(key, value string)) *ChromiumManager {
	if logger == nil {
		logger = func(key, value string) {}
	}

	homeDir, _ := os.UserHomeDir()
	installDir := filepath.Join(homeDir, ".local", "share", config.InstallDirName)

	return &ChromiumManager{
		installDir: installDir,
		logger:     logger,
		config:     config,
	}
}
