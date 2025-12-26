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

	// Dependencies yang akan didownload (.deb packages)
	Dependencies []Dependency
}

// Dependency berisi info dependency yang perlu didownload
type Dependency struct {
	Name        string // Nama package (untuk logging)
	DebianURL   string // URL untuk Debian-based
	UbuntuURL   string // URL untuk Ubuntu-based (fallback ke DebianURL jika kosong)
	LibraryName string // Nama library file (untuk checking)
}

// DefaultConfig mengembalikan konfigurasi default untuk Chrome for Testing
func DefaultConfig() Config {
	return Config{
		// Chrome for Testing - sudah include banyak dependencies
		DownloadURL:    "https://storage.googleapis.com/chrome-for-testing-public/131.0.6778.204/linux64/chrome-linux64.zip",
		InstallDirName: "chrome-for-testing",
		Version:        "131.0.6778.204",
		Dependencies: []Dependency{
			// Dependencies umum yang dibutuhkan Chrome
			{
				Name:        "libnss3",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/n/nss/libnss3_3.87-1_amd64.deb",
				UbuntuURL:   "http://security.ubuntu.com/ubuntu/pool/main/n/nss/libnss3_3.98-0ubuntu0.22.04.2_amd64.deb",
				LibraryName: "libnss3.so",
			},
			{
				Name:        "libnspr4",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/n/nspr/libnspr4_4.35-1_amd64.deb",
				UbuntuURL:   "http://security.ubuntu.com/ubuntu/pool/main/n/nspr/libnspr4_4.32-3ubuntu0.22.04.1_amd64.deb",
				LibraryName: "libnspr4.so",
			},
			{
				Name:        "libatk-bridge2.0-0",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/a/at-spi2-core/libatk-bridge2.0-0_2.46.0-5_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/a/at-spi2-core/libatk-bridge2.0-0_2.38.0-3_amd64.deb",
				LibraryName: "libatk-bridge-2.0.so",
			},
			{
				Name:        "libatk1.0-0",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/a/atk1.0/libatk1.0-0_2.46.0-5_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/a/atk1.0/libatk1.0-0_2.36.0-3build1_amd64.deb",
				LibraryName: "libatk-1.0.so",
			},
			{
				Name:        "libatspi2.0-0",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/a/at-spi2-core/libatspi2.0-0_2.46.0-5_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/a/at-spi2-core/libatspi2.0-0_2.38.0-3_amd64.deb",
				LibraryName: "libatspi.so",
			},
			{
				Name:        "libcups2",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/c/cups/libcups2_2.4.2-3_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/c/cups/libcups2_2.4.1-1ubuntu4.4_amd64.deb",
				LibraryName: "libcups.so",
			},
			{
				Name:        "libdrm2",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/libd/libdrm/libdrm2_2.4.114-1_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/libd/libdrm/libdrm2_2.4.113-2~ubuntu0.22.04.1_amd64.deb",
				LibraryName: "libdrm.so",
			},
			{
				Name:        "libgbm1",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/m/mesa/libgbm1_22.3.6-1_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/m/mesa/libgbm1_23.2.1-1ubuntu3.1~22.04.2_amd64.deb",
				LibraryName: "libgbm.so",
			},
			{
				Name:        "libasound2",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/a/alsa-lib/libasound2_1.2.8-1_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/a/alsa-lib/libasound2_1.2.6.1-1ubuntu1.1_amd64.deb",
				LibraryName: "libasound.so",
			},
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
