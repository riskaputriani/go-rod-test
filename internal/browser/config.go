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
			// Dependencies umum yang dibutuhkan Chrome (Debian 11 Bullseye versions)
			{
				Name:        "libnss3",
				DebianURL:   "https://ftp.debian.org/debian/pool/main/n/nss/libnss3_3.110-1_amd64.deb",
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
				DebianURL:   "http://ftp.debian.org/debian/pool/main/a/at-spi2-atk/libatk-bridge2.0-0_2.38.0-1_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/a/at-spi2-core/libatk-bridge2.0-0_2.38.0-3_amd64.deb",
				LibraryName: "libatk-bridge-2.0.so",
			},
			{
				Name:        "libatk1.0-0",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/a/atk1.0/libatk1.0-0_2.36.0-2_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/a/atk1.0/libatk1.0-0_2.36.0-3build1_amd64.deb",
				LibraryName: "libatk-1.0.so",
			},
			{
				Name:        "libatspi2.0-0",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/a/at-spi2-core/libatspi2.0-0_2.38.0-4+deb11u1_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/a/at-spi2-core/libatspi2.0-0_2.38.0-3_amd64.deb",
				LibraryName: "libatspi.so",
			},
			{
				Name:        "libcups2",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/c/cups/libcups2_2.3.3op2-3+deb11u8_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/c/cups/libcups2_2.4.1-1ubuntu4.4_amd64.deb",
				LibraryName: "libcups.so",
			},
			{
				Name:        "libdrm2",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/libd/libdrm/libdrm2_2.4.104-1_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/libd/libdrm/libdrm2_2.4.113-2~ubuntu0.22.04.1_amd64.deb",
				LibraryName: "libdrm.so",
			},
			{
				Name:        "libgbm1",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/m/mesa/libgbm1_20.3.5-1_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/m/mesa/libgbm1_23.2.1-1ubuntu3.1~22.04.2_amd64.deb",
				LibraryName: "libgbm.so",
			},
			{
				Name:        "libasound2",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/a/alsa-lib/libasound2_1.2.4-1.1_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/a/alsa-lib/libasound2_1.2.6.1-1ubuntu1.1_amd64.deb",
				LibraryName: "libasound.so",
			},
			// X11 Libraries yang dibutuhkan Chrome
			{
				Name:        "libxcomposite1",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/libx/libxcomposite/libxcomposite1_0.4.5-1_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/libx/libxcomposite/libxcomposite1_0.4.5-1_amd64.deb",
				LibraryName: "libXcomposite.so",
			},
			{
				Name:        "libxdamage1",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/libx/libxdamage/libxdamage1_1.1.5-2_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/libx/libxdamage/libxdamage1_1.1.5-2_amd64.deb",
				LibraryName: "libXdamage.so",
			},
			{
				Name:        "libxext6",
				DebianURL:   "https://ftp.debian.org/debian/pool/main/libx/libxext/libxext6_1.3.3-1.1_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/libx/libxext/libxext6_1.3.4-1build1_amd64.deb",
				LibraryName: "libXext.so",
			},
			{
				Name:        "libxfixes3",
				DebianURL:   "https://ftp.debian.org/debian/pool/main/libx/libxfixes/libxfixes3_6.0.0-2_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/libx/libxfixes/libxfixes3_6.0.0-1_amd64.deb",
				LibraryName: "libXfixes.so",
			},
			{
				Name:        "libxrandr2",
				DebianURL:   "https://ftp.debian.org/debian/pool/main/libx/libxrandr/libxrandr2_1.5.1-1_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/libx/libxrandr/libxrandr2_1.5.2-2build1_amd64.deb",
				LibraryName: "libXrandr.so",
			},
			{
				Name:        "libxtst6",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/libx/libxtst/libxtst6_1.2.3-1_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/libx/libxtst/libxtst6_1.2.3-1build4_amd64.deb",
				LibraryName: "libXtst.so",
			},
			{
				Name:        "libxss1",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/libx/libxss/libxss1_1.2.3-1_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/libx/libxss/libxss1_1.2.3-1build2_amd64.deb",
				LibraryName: "libXss.so",
			},
			{
				Name:        "libcairo2",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/c/cairo/libcairo2_1.16.0-5_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/c/cairo/libcairo2_1.16.0-5ubuntu2_amd64.deb",
				LibraryName: "libcairo.so",
			},
			{
				Name:        "libpango-1.0-0",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/p/pango1.0/libpango-1.0-0_1.46.2-3_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/p/pango1.0/libpango-1.0-0_1.50.6+ds-2_amd64.deb",
				LibraryName: "libpango-1.0.so",
			},
			{
				Name:        "libglib2.0-0",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/g/glib2.0/libglib2.0-0_2.66.8-1+deb11u4_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/g/glib2.0/libglib2.0-0_2.72.1-1_amd64.deb",
				LibraryName: "libglib-2.0.so",
			},
			{
				Name:        "libavahi-common3",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/a/avahi/libavahi-common3_0.8-5+deb11u2_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/a/avahi/libavahi-common3_0.8-5ubuntu5_amd64.deb",
				LibraryName: "libavahi-common.so",
			},
			{
				Name:        "libavahi-client3",
				DebianURL:   "http://ftp.debian.org/debian/pool/main/a/avahi/libavahi-client3_0.8-5+deb11u2_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/a/avahi/libavahi-client3_0.8-5ubuntu5_amd64.deb",
				LibraryName: "libavahi-client.so",
			},
			{
				Name:        "libpcre3",
				DebianURL:   "https://ftp.debian.org/debian/pool/main/p/pcre3/libpcre3_8.39-13_amd64.deb",
				UbuntuURL:   "http://archive.ubuntu.com/ubuntu/pool/main/p/pcre3/libpcre3_2%3a8.39-13ubuntu0.22.04.1_amd64.deb",
				LibraryName: "libpcre.so",
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
