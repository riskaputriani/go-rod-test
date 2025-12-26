package browser

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
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

// extractTarXz mengekstrak file tar.xz
func (cm *ChromiumManager) extractTarXz(r io.Reader) error {
	// Karena Go tidak memiliki dukungan bawaan untuk XZ,
	// kita akan menggunakan command eksternal jika tersedia
	// Atau kita bisa save file dulu lalu ekstrak

	// Simpan file sementara
	tmpFile := filepath.Join(os.TempDir(), "chromium.tar.xz")
	f, err := os.Create(tmpFile)
	if err != nil {
		return err
	}

	cm.logger("chromium_temp_file", tmpFile)

	// Copy data ke file
	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		os.Remove(tmpFile)
		return err
	}
	f.Close()

	// Ekstrak menggunakan command line
	// Untuk Linux, gunakan tar command
	if runtime.GOOS == "linux" {
		return cm.extractUsingTarCommand(tmpFile)
	}

	// Untuk OS lain, return error
	return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
}

// extractUsingTarCommand mengekstrak menggunakan tar command
func (cm *ChromiumManager) extractUsingTarCommand(tarFile string) error {
	cm.logger("chromium_extract_cmd", fmt.Sprintf("tar -xf %s -C %s", tarFile, cm.installDir))

	// Gunakan exec untuk menjalankan tar
	cmd := exec.Command("tar", "-xf", tarFile, "-C", cm.installDir)
	output, err := cmd.CombinedOutput()

	if err != nil {
		cm.logger("chromium_extract_error", string(output))
		return fmt.Errorf("failed to extract: %w (output: %s)", err, string(output))
	}

	// Hapus file temporary setelah ekstrak
	os.Remove(tarFile)
	cm.logger("chromium_extract", "success")

	return nil
}

// GetBrowser membuat dan mengembalikan instance browser Rod
func (cm *ChromiumManager) GetBrowser() (*rod.Browser, error) {
	// Setup Chromium jika belum
	if err := cm.Setup(); err != nil {
		return nil, err
	}

	// Jika execPath tidak ditemukan, gunakan launcher default
	var u string

	if cm.execPath != "" {
		// Gunakan Ungoogled Chromium
		cm.logger("browser_using", "ungoogled_chromium")
		u = launcher.New().
			Bin(cm.execPath).
			Headless(true).
			NoSandbox(true).
			MustLaunch()
	} else {
		// Fallback ke browser default
		cm.logger("browser_using", "default")
		path, _ := launcher.LookPath()
		u = launcher.New().
			Bin(path).
			Headless(true).
			NoSandbox(true).
			MustLaunch()
	}

	// Buat browser
	browser := rod.New().ControlURL(u).MustConnect()

	cm.logger("browser_status", "connected")
	return browser, nil
}

// GetExecutablePath mengembalikan path ke executable Chromium
func (cm *ChromiumManager) GetExecutablePath() string {
	return cm.execPath
}
