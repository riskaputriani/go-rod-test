package browser

import (
	"os"
	"path/filepath"
)

// Config berisi konfigurasi untuk Ungoogled Chromium
type Config struct {
	// URL untuk download Ungoogled Chromium
	DownloadURL string

	// Nama direktori instalasi
	InstallDirName string

	// Versi Chromium
	Version string
}

// DefaultConfig mengembalikan konfigurasi default
func DefaultConfig() Config {
	return Config{
		DownloadURL:    "https://github.com/macchrome/linchrome/releases/download/v142.0.7444.229-M142.0.7444.229-r1522585-portable-ungoogled-Lin64/ungoogled-chromium_142.0.7444.229_1.vaapi_linux.tar.xz",
		InstallDirName: "ungoogled-chromium",
		Version:        "142.0.7444.229",
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
