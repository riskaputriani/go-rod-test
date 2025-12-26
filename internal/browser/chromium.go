package browser

import (
	"archive/tar"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/ulikunitz/xz"
)

// ChromiumManager mengelola instalasi dan konfigurasi Ungoogled Chromium
type ChromiumManager struct {
	installDir string
	execPath   string
	logger     func(key, value string)
	config     Config
}

// NewChromiumManager membuat instance baru ChromiumManager
func NewChromiumManager(logger func(key, value string)) *ChromiumManager {
	return NewChromiumManagerWithConfig(DefaultConfig(), logger)
}

// Setup mengecek dan mengunduh Ungoogled Chromium jika belum ada
func (cm *ChromiumManager) Setup() error {
	cm.logger("chromium_install_dir", cm.installDir)

	// Cek apakah sudah terinstall
	if cm.isInstalled() {
		cm.logger("chromium_status", "already_installed")
		return nil
	}

	cm.logger("chromium_status", "not_found_downloading")

	// Download dan ekstrak
	if err := cm.downloadAndExtract(); err != nil {
		return fmt.Errorf("failed to setup chromium: %w", err)
	}

	cm.logger("chromium_status", "installation_complete")
	return nil
}

// isInstalled mengecek apakah Chromium sudah terinstall
func (cm *ChromiumManager) isInstalled() bool {
	// Cari executable chrome di direktori instalasi
	chromePath := filepath.Join(cm.installDir, "chrome")

	if _, err := os.Stat(chromePath); err == nil {
		cm.execPath = chromePath
		cm.logger("chromium_executable", chromePath)
		return true
	}

	// Cari di subdirektori (karena hasil ekstrak mungkin punya folder tambahan)
	entries, err := os.ReadDir(cm.installDir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			chromePath := filepath.Join(cm.installDir, entry.Name(), "chrome")
			if _, err := os.Stat(chromePath); err == nil {
				cm.execPath = chromePath
				cm.logger("chromium_executable", chromePath)
				return true
			}
		}
	}

	return false
}

// downloadAndExtract mengunduh dan mengekstrak Ungoogled Chromium
func (cm *ChromiumManager) downloadAndExtract() error {
	// Buat direktori instalasi
	if err := os.MkdirAll(cm.installDir, 0o755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	// Download file
	cm.logger("chromium_download", "starting")
	cm.logger("chromium_url", cm.config.DownloadURL)
	resp, err := http.Get(cm.config.DownloadURL)
	if err != nil {
		return fmt.Errorf("failed to download chromium: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download chromium: status %d", resp.StatusCode)
	}

	cm.logger("chromium_download", "extracting")

	// Ekstrak langsung ke direktori instalasi
	if err := cm.extractTarXz(resp.Body); err != nil {
		return fmt.Errorf("failed to extract chromium: %w", err)
	}

	// Verifikasi instalasi
	if !cm.isInstalled() {
		return fmt.Errorf("chromium executable not found after extraction")
	}

	return nil
}

// extractTarXz mengekstrak file tar.xz menggunakan pure Go (tanpa dependency eksternal)
func (cm *ChromiumManager) extractTarXz(r io.Reader) error {
	cm.logger("chromium_extract", "using_pure_go_no_external_dependency")

	// Decompress xz
	xzReader, err := xz.NewReader(r)
	if err != nil {
		return fmt.Errorf("failed to create xz reader: %w", err)
	}

	// Extract tar
	tarReader := tar.NewReader(xzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar: %w", err)
		}

		target := filepath.Join(cm.installDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory
			if err := os.MkdirAll(target, 0o755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", target, err)
			}
			cm.logger("chromium_extract_dir", header.Name)

		case tar.TypeReg:
			// Create file
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", target, err)
			}

			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", target, err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write file %s: %w", target, err)
			}

			outFile.Close()

			// Log progress setiap 100 file
			if header.Name[len(header.Name)-1] == 'e' {
				cm.logger("chromium_extract_file", header.Name)
			}

		case tar.TypeSymlink:
			// Create symlink
			if err := os.Symlink(header.Linkname, target); err != nil {
				// Ignore jika symlink sudah ada
				if !os.IsExist(err) {
					return fmt.Errorf("failed to create symlink %s: %w", target, err)
				}
			}
		}
	}

	cm.logger("chromium_extract", "success")
	return nil
}

// GetBrowser membuat dan mengembalikan instance browser Rod
func (cm *ChromiumManager) GetBrowser() (*rod.Browser, error) {
	// Setup Chromium jika belum (akan download jika perlu)
	if err := cm.Setup(); err != nil {
		return nil, fmt.Errorf("failed to setup chromium: %w", err)
	}

	// Pastikan execPath sudah ditemukan
	if cm.execPath == "" {
		return nil, fmt.Errorf("chromium executable not found after setup")
	}

	// Gunakan Ungoogled Chromium yang sudah didownload
	cm.logger("browser_using", "ungoogled_chromium")
	cm.logger("browser_executable", cm.execPath)

	u := launcher.New().
		Bin(cm.execPath).
		Headless(true).
		NoSandbox(true).
		MustLaunch()

	// Buat browser
	browser := rod.New().ControlURL(u).MustConnect()

	cm.logger("browser_status", "connected")
	return browser, nil
}

// GetExecutablePath mengembalikan path ke executable Chromium
func (cm *ChromiumManager) GetExecutablePath() string {
	return cm.execPath
}
